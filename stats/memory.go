package main

import (
	"fmt"
	"my"
	"strconv"
	"strings"
)

// memryInfo : displays RAM stats
//	swap info is hidden until it is in use
func memoryInfo() {
	lines := getFromFile("/proc/meminfo")
	memMap := make(map[string]int)
	for _, line := range lines {
		strArray := strings.Fields(line)
		if len(strArray) > 0 {
			memMap[strArray[0]], _ = strconv.Atoi(string(strArray[1]))
		}
	}

	memused := memMap["MemTotal:"] - memMap["MemFree:"] - memMap["Buffers:"] - memMap["Cached:"]
	percent := float64(memused) / float64(memMap["MemTotal:"]) * 100.0
	// NB. Memories are in k
	fmt.Printf("${color yellow}${alignc}RAM: %5s M ${color cyan}(%5.2f%%) / %s M\n", my.PrintCommaInt(memused/1024), percent, my.PrintCommaInt(memMap["MemTotal:"]/1024))
	printgraph(lineWidth, percent)

	if memMap["SwapTotal:"] != memMap["SwapFree:"] {
		swapused := memMap["SwapTotal:"] - memMap["SwapFree:"]
		percent := float64(swapused) / float64(memMap["SwapTotal:"]) * 100.0
		fmt.Printf(" ${color red}Swap: %5s k ${color cyan}(%5.2f%%) / %s M\n", my.PrintCommaInt(swapused), percent, my.PrintCommaInt(memMap["SwapTotal:"]/1024))
		printgraph(lineWidth, percent)
	}

}
