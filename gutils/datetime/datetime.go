package datetime

import (
	"math"
	"time"
)

// Parse 日期转为时间 "20060102", "20240731"
func Parse(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}

// ToCSTTime 时间转成北京时间
func ToCSTTime(date time.Time) time.Time {
	cstLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return date
	}

	return date.In(cstLocation)
}

// ToUTCTime 时间转成utc时间
func ToUTCTime(date time.Time) time.Time {
	return date.UTC()
}

// GetNowDate 返回当前日期的格式 yyyy-mm-dd
func GetNowDate() string {
	return time.Now().Format(time.DateOnly)
}

// ParseDateToTime 解析日期字符串为时间, 例如"2024-07-19"
func ParseDateToTime(dateStr string) (time.Time, error) {
	return time.Parse(time.DateOnly, dateStr)
}

// GetNowTime 返回当前时间的格式 hh-mm-ss
func GetNowTime() string {
	return time.Now().Format(time.TimeOnly)
}

// GetNowDateTime 返回当前日期时间的格式 yyyy-mm-dd hh-mm-ss
func GetNowDateTime() string {
	return time.Now().Format(time.DateTime)
}

// AddMinute 给指定日期添加或减少分钟
func AddMinute(datetime time.Time, minute int64) time.Time {
	return datetime.Add(time.Minute * time.Duration(minute))
}

// AddHour 给指定日期添加或减少小时
func AddHour(datetime time.Time, hour int64) time.Time {
	return datetime.Add(time.Hour * time.Duration(hour))
}

// AddDay 给指定日期添加或减少天
func AddDay(datetime time.Time, day int64) time.Time {
	return datetime.Add(24 * time.Hour * time.Duration(day))
}

// AddYear 给指定日期添加或减少年
func AddYear(datetime time.Time, year int64) time.Time {
	return datetime.Add(365 * 24 * time.Hour * time.Duration(year))
}

// GetStartOfDay 获取给定日期的一天的开始时间（零点）
func GetStartOfDay(datetime time.Time) time.Time {
	result := time.Date(
		datetime.Year(),
		datetime.Month(),
		datetime.Day(),
		0,
		0,
		0,
		0,
		datetime.Location(),
	)

	return result
}

// GetEndOfDay 获取给定日期的一天的最后时间（23点59分59秒）
func GetEndOfDay(datetime time.Time) time.Time {
	result := time.Date(
		datetime.Year(),
		datetime.Month(),
		datetime.Day(),
		23,
		59,
		59,
		0,
		datetime.Location(),
	)

	return result
}

// GetAnyDateTime 获取任意日期的时间
func GetAnyDateTime(year, month, day, hour, minute, second int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

// GetStartOfDatetime 获取指定日期之前若干天的开始时间（零点）
func GetStartOfDatetime(datetime time.Time, daysAgo int) time.Time {
	startOfDatetime := datetime.AddDate(0, 0, -daysAgo)
	result := time.Date(
		startOfDatetime.Year(),
		startOfDatetime.Month(),
		startOfDatetime.Day(),
		0,
		0,
		0,
		0,
		startOfDatetime.Location(),
	)

	return result
}

// GetEndOfDatetime 获取指定日期之后若干天的结束时间（23:59:59）
func GetEndOfDatetime(datetime time.Time, daysLater int) time.Time {
	endOfDatetime := datetime.AddDate(0, 0, daysLater)
	result := time.Date(
		endOfDatetime.Year(),
		endOfDatetime.Month(),
		endOfDatetime.Day(),
		23,
		59,
		59,
		0,
		endOfDatetime.Location(),
	)

	return result
}

// GetStartOfWeek 获取给定日期所在周的周一时间（零点）
func GetStartOfWeek(datetime time.Time) time.Time {
	weekday := int(datetime.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	daysSinceMonday := weekday - 1
	startOfWeek := datetime.AddDate(0, 0, -daysSinceMonday)

	result := time.Date(
		startOfWeek.Year(),
		startOfWeek.Month(),
		startOfWeek.Day(),
		0,
		0,
		0,
		0,
		startOfWeek.Location(),
	)

	return result
}

// GetEndOfWeek 获取给定日期所在周的周日时间（23点59分59秒）
func GetEndOfWeek(datetime time.Time) time.Time {
	weekday := int(datetime.Weekday())
	daysUntilSunday := (7 - weekday) % 7
	endOfWeek := datetime.AddDate(0, 0, daysUntilSunday)

	result := time.Date(
		endOfWeek.Year(),
		endOfWeek.Month(),
		endOfWeek.Day(),
		23,
		59,
		59,
		0,
		endOfWeek.Location(),
	)

	return result
}

// GetStartOfMonth 获取给定日期所在月的月初时间（即该月的第一天的零点）
func GetStartOfMonth(datetime time.Time) time.Time {
	year, month, _ := datetime.Date()
	result := time.Date(
		year,
		month,
		1,
		0,
		0,
		0,
		0,
		datetime.Location(),
	)

	return result
}

// GetEndOfMonth 获取给定日期所在月的月末时间（即该月的最后一天的23点59分59秒）
func GetEndOfMonth(datetime time.Time) time.Time {
	endOfNextMonth := GetStartOfMonth(datetime).AddDate(0, 1, 0)
	result := endOfNextMonth.Add(-time.Second)

	return result
}

// GetStartOfQuarter 获取给定日期所在季度的季度开始时间（即该季度的第一个月的第一天的零点）
func GetStartOfQuarter(datetime time.Time) time.Time {
	quarterMonth := (int(datetime.Month())-1)/3*3 + 1
	result := time.Date(
		datetime.Year(),
		time.Month(quarterMonth),
		1,
		0,
		0,
		0,
		0,
		datetime.Location(),
	)

	return result
}

// GetEndOfQuarter 获取给定日期所在季度的季度末时间（即该季度的最后一个月的最后一天的23点59分59秒）
func GetEndOfQuarter(datetime time.Time) time.Time {
	quarterMonth := (int(datetime.Month())-1)/3*3 + 1
	startOfNextQuarter := time.Date(
		datetime.Year(),
		time.Month(quarterMonth)+3,
		1,
		0,
		0,
		0,
		0,
		datetime.Location(),
	)
	result := startOfNextQuarter.Add(-time.Second)

	return result
}

// GetStartOfYear 获取给定日期所在年的年初时间（即该年的第一天的零点）
func GetStartOfYear(datetime time.Time) time.Time {
	result := time.Date(
		datetime.Year(),
		time.January,
		1,
		0,
		0,
		0,
		0,
		datetime.Location(),
	)

	return result
}

// GetEndOfYear 获取给定日期所在年的年末时间（即该年的最后一天的23点59分59秒）
func GetEndOfYear(datetime time.Time) time.Time {
	endNextYear := GetStartOfYear(datetime).AddDate(1, 0, 0)
	result := endNextYear.Add(-time.Second)

	return result
}

// IsLeapYear 是否是闰年
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// CurrentSubStartTime 计算时间差, 返回毫秒时间差(两位小数)
func CurrentSubStartTime(startTime time.Time) float64 {
	startTimeN := startTime.UnixNano()
	endTimeN := time.Now().UnixNano()
	return math.Round(float64(endTimeN-startTimeN)/10e6*100) / 100
}
