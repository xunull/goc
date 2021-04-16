package timex

import (
	"fmt"
	"strconv"
	"time"
)

type Escape struct {
	Day    int64
	Hour   int64
	Minute int64
	Second int64
}

func GetEscapeTime(ti time.Duration) Escape {
	escape := int64(ti.Seconds())
	second := escape % 60
	minute := escape / 60 % 60
	hour := escape / (60 * 60) % 24
	day := escape / (60 * 60 * 24)

	return Escape{
		day, hour, minute, second,
	}
}

func GetEscapeTimePlainStr(ti time.Duration) string {
	seconds := int(ti.Seconds())

	second := seconds % 60
	minute := seconds / 60 % 60
	hour := seconds / (60 * 60) % 24

	if hour != 0 {
		return fmt.Sprintf("%d:%.1fh", hour, float64(minute)/float64(60))
	}
	if minute != 0 {
		return fmt.Sprintf("%dm", minute)
	}
	return fmt.Sprintf("%ds", second)

}

func GetTwoTimeEscapeStr(start, end time.Time) string {
	d := end.Sub(start)
	return GetEscapeTimePlainStr(d)
}

func IsZero(t time.Time) bool {
	return t.IsZero()
}

// ---------------------------------------------------------------------------------------------------------------------

func GetYYYYMMDDHHMMSS() string {
	now := time.Now()
	return now.Format("20060102_150405")
}

func GetYYYYMMDD() string {
	return time.Now().Format("20060102")
}

func GetTimeYYYYMMDD(t time.Time) string {
	return t.Format("20060102")
}

// ---------------------------------------------------------------------------------------------------------------------

func GetYYYYMMDDTime(target string) (time.Time, error) {
	return time.Parse("20060102", target)
}

func GetYYYYTime(target string) (time.Time, error) {
	return time.Parse("2006", target)
}

func GetYYYYMMTime(target string) (time.Time, error) {
	return time.Parse("200601", target)
}

func GetTimeByFormat(target, format string) (time.Time, error) {
	return time.Parse(format, target)
}

func GetRFC3339Time(target string) (time.Time, error) {
	return time.Parse(time.RFC3339, target)
}

func GetDay() string {
	return time.Now().Format("20060102")
}

func GetDayInt() (int, error) {
	s := time.Now().Format("20060102")
	return strconv.Atoi(s)
}

func GetDayIntOrPanic() int {
	s := time.Now().Format("20060102")
	res, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return res
}

func GetTomorrowIntOrPanic() int {
	m := time.Now().Add(time.Hour * 24)
	s := m.Format("20060102")
	res, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return res
}

func IsToday(d time.Time) bool {
	return time.Now().Day() == d.Day()
}

func IsBeforeToday(d time.Time) bool {
	t, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	return t.After(d)
}
