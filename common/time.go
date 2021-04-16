package common

import "time"



func GetYYYYMMDD() string {
	return time.Now().Format("20060102")
}

func GetYYYYMMDDHHMMSS() string {
	now := time.Now()
	return now.Format("20060102_150405")
}

func IsToday(d time.Time) bool {
	return time.Now().Day() == d.Day()
}

func IsBeforeToday(d time.Time) bool {
	t, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	return t.After(d)
}
