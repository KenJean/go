package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

//  CPU temperatures
func temperatures() {
	const fnPath1 = "/sys/class/hwmon/hwmon0/temp1_input"
	const fnPath2 = "/sys/class/hwmon/hwmon0/temp2_input"
	// const RaspPath = "/sys/class/thermal/thermal_zone0/temp"
	var onlyOne bool
	var temperature1, temperature2 float64

	content, err := ioutil.ReadFile(fnPath1)
	if err != nil {
		log.Printf("Cannot open %v: %v", fnPath1, err)
	} else {
		num, _ := strconv.Atoi(strings.TrimSpace(string(content)))
		temperature1 = float64(num) / 1000.0

		content, err := ioutil.ReadFile(fnPath2)
		if err != nil {
			log.Printf("Cannot open %v: %v", fnPath2, err)
		}
		onlyOne = false
		num, _ = strconv.Atoi(strings.TrimSpace(string(content)))
		temperature2 = float64(num) / 1000.0
	}
	if onlyOne {
		fmt.Printf("${color yellow}${alignc}CPU: ${color red}%.1f°C\n", temperature1)
	} else {
		fmt.Printf("${color yellow}${alignc}CPU 1: ${color red}%.1f°C${color yellow}    CPU 2: ${color red}%.1f°C\n", temperature1, temperature2)
	}
}
