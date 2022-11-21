package main

import "strconv"

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
