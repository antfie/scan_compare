package main

import "time"

func parseVeracodeDate(date string) time.Time {
	parsed, err := time.Parse("2006-01-02 15:04:05 MST", date)

	if err != nil {
		panic(err)
	}

	return parsed
}
