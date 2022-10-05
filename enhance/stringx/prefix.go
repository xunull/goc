package stringx

import (
	"time"
)

func IsStartWithDate(target string) bool {
	l := len(target)
	var err error

	if l >= 4 {

		temp := target[:4]
		_, err = time.Parse("2006", temp)
		if err != nil {
			if l >= 8 {
				temp := target[:8]
				_, err = time.Parse("20060102", temp)
				if err != nil {
					if l >= 9 {
						temp := target[:9]
						_, err = time.Parse("2006-0102", temp)
						if err == nil {
							return true
						}
						_, err = time.Parse("2006_0102", temp)
						if err == nil {
							return true
						}
					}
				} else {
					return true
				}
			}
		} else {
			return true
		}
	} else {
		return false
	}
	return false
}
