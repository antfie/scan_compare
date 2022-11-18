package main

func isStringInStringArray(input string, list []string) bool {
	for _, item := range list {
		if input == item {
			return true
		}
	}

	return false
}
