package main

import (
	"log"

	"golang.org/x/sys/windows/registry"
)

const (
	// https://docs.microsoft.com/en-us/windows-hardware/drivers/kernel/interrupt-affinity-and-priority
	IrqPolicyMachineDefault                    = iota // 0
	IrqPolicyAllCloseProcessors                       // 1
	IrqPolicyOneCloseProcessor                        // 2
	IrqPolicyAllProcessorsInMachine                   // 3
	IrqPolicySpecifiedProcessors                      // 4
	IrqPolicySpreadMessagesAcrossAllProcessors        // 5
)

type Bits uint64

const ZeroBit = Bits(0)

func Set(b, flag Bits) Bits    { return b | flag }
func Clear(b, flag Bits) Bits  { return b &^ flag }
func Toggle(b, flag Bits) Bits { return b ^ flag }
func Has(b, flag Bits) bool    { return b&flag != 0 }

func setAffinityPolicy(item *Device) {
	var k registry.Key
	var err error

	if item.DevicePolicy == 0 {

		k, err = registry.OpenKey(item.reg, `Interrupt Management\Affinity Policy`, registry.ALL_ACCESS)
		if err != nil {
			log.Println(err)
		}

		if err := registry.DeleteKey(item.reg, `Interrupt Management\Affinity Policy`); err != nil {
			log.Println(err)
		}

	} else {

		k, _, err = registry.CreateKey(item.reg, `Interrupt Management\Affinity Policy`, registry.ALL_ACCESS)
		if err != nil {
			log.Println(err)
		}

		if err := k.SetDWordValue("DevicePolicy", item.DevicePolicy); err != nil {
			log.Println(err)
		}

		if item.DevicePolicy != 4 {
			k.DeleteValue("AssignmentSetOverride")
		}

		AssignmentSetOverrideByte := i64tob(uint64(item.AssignmentSetOverride))
		if err := k.SetBinaryValue("AssignmentSetOverride", AssignmentSetOverrideByte[:clen(AssignmentSetOverrideByte)]); err != nil {
			log.Println(err)
		}

	}
	if err := k.Close(); err != nil {
		log.Println(err)
	}
}

func clen(n []byte) int {
	for i := len(n) - 1; i >= 0; i-- {
		if n[i] != 0 {
			return i + 1
		}
	}
	return len(n)
}

// https://gist.github.com/chiro-hiro/2674626cebbcb5a676355b7aaac4972d
func i64tob(val uint64) []byte {
	r := make([]byte, 8)
	for i := uint64(0); i < 8; i++ {
		r[i] = byte((val >> (i * 8)) & 0xff)
	}
	return r
}

func btoi64(val []byte) uint64 {
	r := uint64(0)
	for i := uint64(0); i < 8; i++ {
		r |= uint64(val[i]) << (8 * i)
	}
	return r
}

func btoi32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 4; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}

func GetDWORDuint32Value(key registry.Key, name string) uint32 {
	buf := make([]byte, 4)
	key.GetValue(name, buf)
	return btoi32(buf)
}

func GetBinaryValue(key registry.Key, name string) []byte {
	value, _, err := key.GetBinaryValue(name)
	if err != nil {
		return []byte{}
	}
	return value
}

func SetRunOnce(val string) {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\RunOnce`, registry.ALL_ACCESS)
	if err != nil {
		log.Println(err)
	}

	err = k.SetStringValue("AutoGpuAffinity", val)
	if err != nil {
		panic(err)
	}

	if err := k.Close(); err != nil {
		log.Println(err)
	}
}
