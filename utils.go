package main

import (
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func colorPrintf(format string) {
	color.New().Printf(format)
}

func isStringInStringArray(input string, list []string) bool {
	for _, item := range list {
		if input == item {
			return true
		}
	}

	return false
}

func stringToFloat(input string) (float64, error) {
	return strconv.ParseFloat(input, 64)
}

func getSortedIntArrayAsSFormattedString(list []int) string {
	sort.Ints(list[:])
	var output []string
	for _, x := range list {
		output = append(output, strconv.Itoa(x))
	}

	return strings.Join(output, ",")
}

func isInIntArray(x int, y []int) bool {
	for _, z := range y {
		if x == z {
			return true
		}
	}

	return false
}

func getFormattedOnlyInSideString(side string) string {
	if side == "A" {
		return color.HiGreenString("Only in A")
	}

	return color.HiMagentaString("Only in B")
}

func getFormattedSideString(side string) string {
	if side == "A" {
		return color.HiGreenString("A")
	}

	return color.HiMagentaString("B")
}

func getFormattedSideStringWithMessage(side, message string) string {
	if side == "A" {
		return color.HiGreenString(message)
	}

	return color.HiMagentaString(message)
}
