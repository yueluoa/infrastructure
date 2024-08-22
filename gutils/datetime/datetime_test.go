package datetime

import (
	"fmt"
	"testing"
	"time"
)

func TestDatetime(t *testing.T) {
	result1 := GetStartOfDay(time.Now())
	fmt.Println(result1)

	result2 := GetEndOfDay(time.Now())
	fmt.Println(result2)

	result3 := GetStartOfWeek(time.Now())
	fmt.Println(result3)

	result4 := GetEndOfWeek(time.Now())
	fmt.Println(result4)

	result5 := GetStartOfMonth(time.Now())
	fmt.Println(result5)

	result6 := GetEndOfMonth(time.Now())
	fmt.Println(result6)

	result7 := GetStartOfQuarter(time.Now())
	fmt.Println(result7)

	result8 := GetEndOfQuarter(time.Now())
	fmt.Println(result8)

	result9 := GetStartOfYear(time.Now())
	fmt.Println(result9)

	result10 := GetEndOfYear(time.Now())
	fmt.Println(result10)

	result11 := AddMinute(time.Now(), 10)
	fmt.Println(result11)
	result11 = AddMinute(time.Now(), -10)
	fmt.Println(result11)

	result12 := AddHour(time.Now(), 10)
	fmt.Println(result12)
	result12 = AddHour(time.Now(), -10)
	fmt.Println(result12)

	result13 := AddDay(time.Now(), 10)
	fmt.Println(result13)
	result13 = AddDay(time.Now(), -10)
	fmt.Println(result13)

	result14 := AddYear(time.Now(), 10)
	fmt.Println(result14)
	result14 = AddYear(time.Now(), -10)
	fmt.Println(result14)

	fmt.Println(IsLeapYear(2025))

	result15 := GetStartOfDatetime(time.Now(), 2)
	fmt.Println(result15)
	result15 = GetEndOfDatetime(time.Now(), 2)
	fmt.Println(result15)

	year := 2024
	month := 7
	day := 15
	hour := 14
	minute := 30
	second := 0
	result16 := GetAnyDateTime(year, month, day, hour, minute, second)
	fmt.Println(result16)
}

// CurrentSubStartTime 计算时间差, 返回秒时间差值(两位小数)
func Test_CurrentSubStartTime(t *testing.T) {
	startTime := time.Now()
	time.Sleep(1*time.Second + 220*time.Millisecond)
	duration := CurrentSubStartTime(startTime)
	fmt.Printf("%.2fms", duration)
}
