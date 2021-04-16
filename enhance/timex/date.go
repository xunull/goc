package timex

import (
	"strconv"
	"time"
)

func GetDayListFromIntDayList(list []int) []string {
	res := make([]string, 0, len(list))
	for _, item := range list {
		t := strconv.Itoa(item)
		res = append(res, strconv.Itoa(item)[len(t)-2:len(t)])
	}
	return res
}

func GetTwoDateDayIntList(start, end time.Time) ([]int, error) {

	leftDay, err := time.Parse("20060102", start.Format("20060102"))
	if err != nil {
		return nil, err
	}
	rightDay, err := time.Parse("20060102", end.Format("20060102"))
	if err != nil {
		return nil, err
	}
	res := make([]int, 0)
	for cur := leftDay; cur.Before(rightDay.Add(24 * time.Hour)); {
		day, err := strconv.Atoi(cur.Format("20060102"))
		if err != nil {
			return nil, err
		}
		res = append(res, day)
		cur = cur.Add(24 * time.Hour)
	}
	return res, nil
}
