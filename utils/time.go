package utils

import "time"

func EndOfThisMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC)
}
