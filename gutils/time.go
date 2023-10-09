package gutils

import (
	"fmt"
	"time"
)

// 获取当日23点59分59秒的时间
func GetTodayEndTime() time.Time {
	t := time.Now()
	endTime := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())

	return endTime
}

// 获取两天前凌晨时间
func GetTwoDaysAgoMidnight() time.Time {
	now := time.Now()
	t := now.AddDate(0, 0, -2)
	twoDaysAgoMidnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return twoDaysAgoMidnight
}

// 当日凌晨
func YesterdayBeforeDawn() time.Time {
	t := time.Now().AddDate(0, 0, -1)
	yesterdayTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return yesterdayTime
}

func ResolveTime(second int) string {
	var (
		secondsPerMinute = 60
		secondsPerHour   = secondsPerMinute * 60
		secondsPerDay    = secondsPerHour * 24
	)

	day := second / secondsPerDay
	if day > 0 {
		second = second - (day * secondsPerDay)
	}
	hour := second / secondsPerHour
	if hour > 0 {
		second = second - (hour * secondsPerHour)
	}
	minute := second / secondsPerMinute
	if minute > 0 {
		second = second - (minute * secondsPerMinute)
	}

	var strTime string
	if day > 0 {
		strTime = fmt.Sprintf("%v天", day)
	} else {
		if hour > 0 {
			strTime = strTime + fmt.Sprintf("%v小时", hour)
		}
		if minute > 0 {
			strTime = strTime + fmt.Sprintf("%v分", minute)
		}
		if second > 0 && strTime == "" {
			strTime = strTime + fmt.Sprintf("%v秒", second)
		}
	}

	return strTime
}
