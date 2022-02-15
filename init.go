package main

import (
	"log"
	"strconv"
)

var (
	HT       bool
	CPUMap   map[Bits]string
	CPUArray []string
	CPUBits  []Bits

	Phys    int
	Cores   int
	Threads int

	defaultSettings DefaultSettings
	result          Result
	GPUdevices      []Device
	handle          DevInfo
	lava            Program

	tempFolder  string
	trialtime   int
	totaltrials int
	cliMode     bool
)

type DefaultSettings struct {
	GPUdevices []Device
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	Phys, Cores, Threads, _ = computeCoresAndProcessors()
	if Threads > Cores {
		HT = true
	}

	CPUMap = make(map[Bits]string, Threads)
	var index Bits = 1
	for i := 0; i < Threads; i++ {
		indexString := strconv.Itoa(i)
		CPUMap[index] = indexString
		CPUArray = append(CPUArray, indexString)
		CPUBits = append(CPUBits, index)
		index *= 2
	}
}
