package main

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
)

const paddingRight = 2

var header = []string{
	"Name",
	"Max",
	"Avg",
	"Min",

	"1%ile",
	"0.1%ile",
	"0.01%ile",
	"0.005%ile",

	"1% low",
	"0.1% low",
	"0.01% low",
	"0.005% low",
}

var (
	highest = []lipgloss.Color{
		lipgloss.Color("#008000"),
		lipgloss.Color("#00bf00"),
		lipgloss.Color("#00FF00"), // Highest
	}
	// lowest = []lipgloss.Color{
	// 	lipgloss.Color("#ff0000"), // lowest
	// 	lipgloss.Color("#ff4040"),
	// 	lipgloss.Color("#ff8080"),
	// }
	defaultStyle = lipgloss.NewStyle()
)

var CPUPoints map[int]uint
var topCPU1 uint
var topCPU2 uint
var topCPU3 uint

type format struct {
	Array  []float64
	Higher []float64
	Lower  []float64
	Length int
}

type TableSettings struct {
	Data    []format
	Lengths []int
}

func tableOutput(cpuList []CPUResultList) {
	var tableSettings TableSettings
	tableSettings.Lengths = make([]int, len(header))
	tableSettings.Data = make([]format, len(header))
	CPUPoints = make(map[int]uint, len(header))

	for i := 0; i < len(header); i++ {
		tableSettings.Lengths[i] = lipgloss.Width(header[i]) + paddingRight
	}

	for _, v := range cpuList {
		tableSettings.Data[1].Array = append(tableSettings.Data[1].Array, v.Max)
		maxl := lipgloss.Width(fmt.Sprintf("%.2f", v.Max)) + paddingRight
		if maxl > tableSettings.Lengths[1] {
			tableSettings.Lengths[1] = maxl
		}

		tableSettings.Data[2].Array = append(tableSettings.Data[2].Array, v.Avg)
		avgl := lipgloss.Width(fmt.Sprintf("%.2f", v.Avg)) + paddingRight
		if avgl > tableSettings.Lengths[2] {
			tableSettings.Lengths[2] = avgl
		}

		tableSettings.Data[3].Array = append(tableSettings.Data[3].Array, v.Min)
		minl := lipgloss.Width(fmt.Sprintf("%.2f", v.Min)) + paddingRight
		if minl > tableSettings.Lengths[3] {
			tableSettings.Lengths[3] = minl
		}

		tableSettings.Data[4].Array = append(tableSettings.Data[4].Array, v.Percent1)
		percent1l := lipgloss.Width(fmt.Sprintf("%.2f", v.Percent1)) + paddingRight
		if percent1l > tableSettings.Lengths[4] {
			tableSettings.Lengths[4] = percent1l
		}

		tableSettings.Data[5].Array = append(tableSettings.Data[5].Array, v.Percent01)
		percent01l := lipgloss.Width(fmt.Sprintf("%.2f", v.Percent01)) + paddingRight
		if percent01l > tableSettings.Lengths[5] {
			tableSettings.Lengths[5] = percent01l
		}

		tableSettings.Data[6].Array = append(tableSettings.Data[6].Array, v.Percent001)
		percent001l := lipgloss.Width(fmt.Sprintf("%.2f", v.Percent001)) + paddingRight
		if percent001l > tableSettings.Lengths[6] {
			tableSettings.Lengths[6] = percent001l
		}

		tableSettings.Data[7].Array = append(tableSettings.Data[7].Array, v.Percent0005)
		percent0005l := lipgloss.Width(fmt.Sprintf("%.2f", v.Percent0005)) + paddingRight
		if percent0005l > tableSettings.Lengths[7] {
			tableSettings.Lengths[7] = percent0005l
		}

		tableSettings.Data[8].Array = append(tableSettings.Data[8].Array, v.LowsOne)
		lowsOnel := lipgloss.Width(fmt.Sprintf("%.2f", v.LowsOne)) + paddingRight
		if lowsOnel > tableSettings.Lengths[8] {
			tableSettings.Lengths[8] = lowsOnel
		}

		tableSettings.Data[9].Array = append(tableSettings.Data[9].Array, v.LowsPoint1)
		lowsPoint1l := lipgloss.Width(fmt.Sprintf("%.2f", v.LowsPoint1)) + paddingRight
		if lowsPoint1l > tableSettings.Lengths[9] {
			tableSettings.Lengths[9] = lowsPoint1l
		}

		tableSettings.Data[10].Array = append(tableSettings.Data[10].Array, v.LowsPoint01)
		lowsPoint01l := lipgloss.Width(fmt.Sprintf("%.2f", v.LowsPoint01)) + paddingRight
		if lowsPoint01l > tableSettings.Lengths[10] {
			tableSettings.Lengths[10] = lowsPoint01l
		}

		tableSettings.Data[11].Array = append(tableSettings.Data[11].Array, v.LowsPoint005)
		lowsPoint005l := lipgloss.Width(fmt.Sprintf("%.2f", v.LowsPoint005)) + paddingRight
		if lowsPoint005l > tableSettings.Lengths[11] {
			tableSettings.Lengths[11] = lowsPoint005l
		}
	}

	for i := 1; i < len(tableSettings.Data); i++ {
		sort.Float64s(tableSettings.Data[i].Array)
		tableSettings.Data[i].Higher = tableSettings.Data[i].Array[len(tableSettings.Data[i].Array)-3:]
		tableSettings.Data[i].Lower = tableSettings.Data[i].Array[:3]
	}

	for i := 0; i < len(header); i++ {
		fmt.Print(defaultStyle.Bold(true).Underline(true).Width(tableSettings.Lengths[i]).Render(header[i]))
	}
	fmt.Print("\n")

	for i, cpu := range cpuList {
		fmt.Print(defaultStyle.Bold(true).Width(tableSettings.Lengths[0]).Render(fmt.Sprintf("CPU%d", i)))
		tableSettings.GetColor_Max(cpu.Max, 1, i)
		tableSettings.GetColor_Max(cpu.Avg, 2, i)
		tableSettings.GetColor_Max(cpu.Min, 3, i)
		tableSettings.GetColor_Max(cpu.Percent1, 4, i)
		tableSettings.GetColor_Max(cpu.Percent01, 5, i)
		tableSettings.GetColor_Max(cpu.Percent001, 6, i)
		tableSettings.GetColor_Max(cpu.Percent0005, 7, i)
		tableSettings.GetColor_Max(cpu.LowsOne, 8, i)
		tableSettings.GetColor_Max(cpu.LowsPoint1, 9, i)
		tableSettings.GetColor_Max(cpu.LowsPoint01, 10, i)
		tableSettings.GetColor_Max(cpu.LowsPoint005, 11, i)
		fmt.Print("\n")
	}

	for _, points := range CPUPoints {
		if points > topCPU1 {
			topCPU1 = points
		}
	}

	for _, points := range CPUPoints {
		if points > topCPU2 && points != topCPU1 {
			topCPU2 = points
		}
	}

	for _, points := range CPUPoints {
		if points > topCPU3 && points < topCPU2 {
			topCPU3 = points
		}
	}

}

func (tableSettings *TableSettings) GetColor_Max(value float64, i, cpuID int) {
	switch value {
	case tableSettings.Data[i].Higher[0]:
		fmt.Print(defaultStyle.Width(tableSettings.Lengths[i]).Foreground(highest[0]).Render(fmt.Sprintf("%.2f", value)))
		CPUPoints[cpuID] += 1
	case tableSettings.Data[i].Higher[1]:
		fmt.Print(defaultStyle.Width(tableSettings.Lengths[i]).Foreground(highest[1]).Render(fmt.Sprintf("%.2f", value)))
		CPUPoints[cpuID] += 2
	case tableSettings.Data[i].Higher[2]:
		fmt.Print(defaultStyle.Width(tableSettings.Lengths[i]).Foreground(highest[2]).Render(fmt.Sprintf("%.2f", value)))
		CPUPoints[cpuID] += 3
	// case Lower[0]:
	// 	fmt.Print(defaultStyle.Width(tableSettings.Lengths[i]).Foreground(lowest[0]).Render(fmt.Sprintf("%.2f", value)))
	// case Lower[1]:
	// 	fmt.Print(defaultStyle.Width(tableSettings.Lengths[i]).Foreground(lowest[1]).Render(fmt.Sprintf("%.2f", value)))
	// case Lower[2]:
	// 	fmt.Print(defaultStyle.Width(tableSettings.Lengths[i]).Foreground(lowest[2]).Render(fmt.Sprintf("%.2f", value)))
	default:
		fmt.Print(defaultStyle.Width(tableSettings.Lengths[i]).Render(fmt.Sprintf("%.2f", value)))
	}
}
