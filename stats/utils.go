package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

// getFromFile: Reads a text file and returns it in a line-separated []string
func getFromFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Cannot open %v: %v", filename, err)
	}
	return strings.Split(string(content), "\n")
}

// printgraph: prints a line graph using the unicode █ & ░
// NB. percent is a float64 that is converted internally to int
func printgraph(width int, percent float64) {
	used := width * int(percent) / 100
	free := width - used
	result := "${color DDAA00}" + strings.Repeat("█", used) + strings.Repeat("░", free)
	fmt.Println(result)
}
