package main

import (
	"encoding/binary"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// https://github.com/DataDog/gohai/blob/6b668acb50dd0962b0eb2d6b6ba0cd4b2348d5f0/cpu/cpu_windows.go
// https://docs.microsoft.com/de-de/windows/win32/api/sysinfoapi/nf-sysinfoapi-getlogicalprocessorinformation?redirectedfrom=MSDN

// wProcessorArchitecture is a wrapper for the union found in LP_SYSTEM_INFO
// https://docs.microsoft.com/en-us/windows/win32/api/sysinfoapi/ns-sysinfoapi-system_info
type wProcessorArchitecture struct {
	WProcessorArchitecture uint16
	WReserved              uint16
}

// ProcessorArchitecture is an idiomatic wrapper for wProcessorArchitecture
type ProcessorArchitecture uint16

// Idiomatic values for wProcessorArchitecture
const (
	AMD64   ProcessorArchitecture = 9
	ARM     ProcessorArchitecture = 5
	ARM64   ProcessorArchitecture = 12
	IA64    ProcessorArchitecture = 6
	INTEL   ProcessorArchitecture = 0
	UNKNOWN ProcessorArchitecture = 0xffff
)

// SystemInfo is an idiomatic wrapper for LpSystemInfo
type SystemInfo struct {
	Arch                      ProcessorArchitecture
	PageSize                  uint32
	MinimumApplicationAddress uintptr
	MaximumApplicationAddress uintptr
	ActiveProcessorMask       uint
	NumberOfProcessors        uint32
	ProcessorType             uint32
	AllocationGranularity     uint32
	ProcessorLevel            uint16
	ProcessorRevision         uint16
}

// LpSystemInfo is a wrapper for LPSYSTEM_INFO
// https://docs.microsoft.com/en-us/windows/win32/api/sysinfoapi/ns-sysinfoapi-system_info
type lpSystemInfo struct {
	Arch                        wProcessorArchitecture
	DwPageSize                  uint32
	LpMinimumApplicationAddress uintptr
	LpMaximumApplicationAddress uintptr
	DwActiveProcessorMask       uint
	DwNumberOfProcessors        uint32
	DwProcessorType             uint32
	DwAllocationGranularity     uint32
	WProcessorLevel             uint16
	WProcessorRevision          uint16
}

var (
	// Library
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	// Functions
	procGetSystemInfo                  = kernel32.NewProc("GetSystemInfo")
	procGetLogicalProcessorInformation = kernel32.NewProc("GetLogicalProcessorInformation")
)

// GetSystemInfo is an idiomatic wrapper for the GetSystemInfo function from sysinfoapi
// https://docs.microsoft.com/en-us/windows/win32/api/sysinfoapi/nf-sysinfoapi-getsysteminfo
func GetSystemInfo() SystemInfo {
	var info lpSystemInfo
	procGetSystemInfo.Call(uintptr(unsafe.Pointer(&info)))
	return SystemInfo{
		Arch:                      ProcessorArchitecture(info.Arch.WProcessorArchitecture),
		PageSize:                  info.DwPageSize,
		MinimumApplicationAddress: info.LpMinimumApplicationAddress,
		MaximumApplicationAddress: info.LpMinimumApplicationAddress,
		ActiveProcessorMask:       info.DwActiveProcessorMask,
		NumberOfProcessors:        info.DwNumberOfProcessors,
		ProcessorType:             info.DwProcessorType,
		AllocationGranularity:     info.DwAllocationGranularity,
		ProcessorLevel:            info.WProcessorLevel,
		ProcessorRevision:         info.WProcessorRevision,
	}
}

const ERROR_INSUFFICIENT_BUFFER syscall.Errno = 122

const RelationProcessorCore = 0
const RelationNumaNode = 1
const RelationCache = 2
const RelationProcessorPackage = 3
const RelationGroup = 4

type SYSTEM_LOGICAL_PROCESSOR_INFORMATION struct {
	ProcessorMask uintptr
	Relationship  int // enum (int)
	// in the Windows header, this is a union of a byte, a DWORD,
	// and a CACHE_DESCRIPTOR structure
	dataunion [16]byte
}

func countBits(num uint64) (count int) {
	count = 0
	for num > 0 {
		if (num & 0x1) == 1 {
			count++
		}
		num >>= 1
	}
	return
}

func byteArrayToProcessorStruct(data []byte) (info SYSTEM_LOGICAL_PROCESSOR_INFORMATION) {
	info.ProcessorMask = uintptr(binary.LittleEndian.Uint64(data))
	info.Relationship = int(binary.LittleEndian.Uint64(data[8:]))
	copy(info.dataunion[0:16], data[16:32])
	return
}

func computeCoresAndProcessors() (phys int, cores int, processors int, err error) {
	var buflen uint32 = 0
	err = syscall.Errno(0)
	// first, figure out how much we need
	status, _, err := procGetLogicalProcessorInformation.Call(uintptr(0), uintptr(unsafe.Pointer(&buflen)))
	if status == 0 {
		if err != ERROR_INSUFFICIENT_BUFFER {
			// only error we're expecing here is insufficient buffer
			// anything else is a failure
			return
		}
	} else {
		// this shouldn't happen. Errno won't be set (because the fuction)
		// succeeded.  So just return something to indicate we've failed
		return 0, 0, 0, syscall.Errno(1)
	}
	buf := make([]byte, buflen)
	status, _, err = procGetLogicalProcessorInformation.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&buflen)))
	if status == 0 {
		return
	}
	// walk through each of the buffers
	var numaNodeCount int32

	// SYSTEM_LOGICAL_PROCESSOR_INFORMATION_SIZE = 32 (64bit)
	for i := 0; uint32(i) < buflen; i += 32 {
		info := byteArrayToProcessorStruct(buf[i : i+32])
		switch info.Relationship {
		case RelationNumaNode:
			numaNodeCount++

		case RelationProcessorCore:
			cores++
			processors += countBits(uint64(info.ProcessorMask))

		case RelationProcessorPackage:
			phys++
		}
	}
	return
}
