package utils

import "time"

// FakeNow is a fake time for testing
func FakeNow(date string) func() time.Time {
	return func() time.Time {
		tm, err := time.ParseInLocation("2006-01-02 15:04:05", date, time.Local)
		if err != nil {
			panic(err)
		}
		return tm
	}
}
