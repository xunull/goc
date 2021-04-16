package timex

import (
	"fmt"
	"strconv"
	"time"
)

func GuessTime(str string) (time.Time, error) {

	if len(str) <= 2 {
		d, err := strconv.Atoi(str)
		if err != nil {
			return time.Time{}, err
		}
		var r string
		today := GetYYYYMMDD()[:6]
		if d < 10 {
			r = fmt.Sprintf("%s0%d", today[:6], d)
			return GetYYYYMMDDTime(r)
		} else if d <= 31 {
			r = fmt.Sprintf("%s%d", today[:6], d)
			return GetYYYYMMDDTime(r)
		} else {
			return GetYYYYMMDDTime(str)
		}
	} else if len(str) <= 4 {

		d, err := strconv.Atoi(str)
		if err != nil {
			return time.Time{}, err
		}
		var r string
		today := GetYYYYMMDD()[:6]
		if d < 999 {
			r = fmt.Sprintf("%s0%d", today[:4], d)
			return GetYYYYMMDDTime(r)
		} else if d <= 1231 {
			r = fmt.Sprintf("%s%d", today[:4], d)
			return GetYYYYMMDDTime(r)
		} else {
			return GetYYYYMMDDTime(str)
		}

	} else {
		return GetYYYYMMDDTime(str)
	}

}
