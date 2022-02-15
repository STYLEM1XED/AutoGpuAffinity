package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Result struct {
	CPUList [][]CPUResultList
}

type CPUResultList struct {
	Max float64
	Avg float64
	Min float64

	LowsPoint005 float64
	LowsPoint01  float64
	LowsPoint1   float64
	LowsOne      float64

	Percent1    float64
	Percent01   float64
	Percent001  float64
	Percent0005 float64
}

func calc(file string) CPUResultList {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(f)
	var benchmarkTime float64
	var frametimes []float64
	var fps []float64
	var index int = -1
	for fileScanner.Scan() {
		text := strings.Split(fileScanner.Text(), ",")
		if index == -1 {
			index = findPos(text)
			if index == -1 {
				log.Fatalln("msBetweenPresents not Found")
			}
			continue
		}

		s, err := strconv.ParseFloat(text[index], 64)
		if err != nil {
			log.Println("err", err)
			continue
		}
		fps = append(fps, 1000/s)
		frametimes = append(frametimes, s)
		benchmarkTime += s
	}
	lastIndex := len(frametimes)

	var result = CPUResultList{}
	sort.Float64s(fps)
	sort.Float64s(frametimes)
	// from small to big

	for _, low := range []float64{
		1,
		0.1,
		0.01,
		0.005,
	} {
		var wall = low / 100 * benchmarkTime
		var currentTotal float64
		for i := len(frametimes) - 1; i >= 0; i-- {
			present := frametimes[i]
			currentTotal += present
			if currentTotal >= wall {
				var fps = 1000 / present
				switch low {
				case 1:
					result.LowsOne = round(fps)
				case 0.1:
					result.LowsPoint1 = round(fps)
				case 0.01:
					result.LowsPoint01 = round(fps)
				case 0.005:
					result.LowsPoint005 = round(fps)
				}
				break
			}
		}
	}

	result.Max = 1000 / frametimes[0]
	result.Avg = 1000 / (benchmarkTime / float64(lastIndex))
	result.Min = 1000 / frametimes[lastIndex-1]

	var huntertProzent = float64(lastIndex) / 100
	result.Percent1 = fps[lastIndex/100]
	result.Percent01 = fps[roundInt(huntertProzent*0.1)]
	result.Percent001 = fps[roundInt(huntertProzent*0.01)]
	result.Percent0005 = fps[roundInt(huntertProzent*0.005)]

	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	f.Close()

	return result
}

func findPos(textArr []string) int {
	for i, val := range textArr {
		if val == "msBetweenPresents" {
			return i
		}
	}
	return -1
}

func roundInt(val float64) int {
	// return int(math.Ceil(val))
	return int(math.Round(val))
}

func round(val float64) float64 {
	return math.Round(val*100) / 100
}
