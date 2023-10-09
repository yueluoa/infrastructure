package html

import (
	"html/template"
	"math"
	"strings"
	"time"
)

// 自定义模板函数
var hashFunc = template.FuncMap{
	"formatTime": formatTime,
	"suffixTime": suffixTime,
	"format":     format,
	"add":        add,
	"mod":        mod,
	"sub":        sub,
	"ride":       ride,
	"divide":     divide,
	"explode":    explode,
	"noescape":   noescape,
}

func formatTime(timeStr string) string {
	t, _ := time.Parse("2006-01-02 15:04:05", timeStr)
	return t.Format(time.RFC3339) + "+08:00"
}

func suffixTime() string {
	return time.Now().Format("15:04")
}

func format(i interface{}, format string) string {
	switch i.(type) {
	case time.Time:
		return (i.(time.Time)).Format(format)
	case int64:
		val := i.(int64)
		return time.Unix(val, 0).Format(format)
	}

	return ""
}

func add(nums ...interface{}) int {
	var total int
	for _, num := range nums {
		if n, ok := num.(int); ok {
			total += n
		}
	}
	return total
}

func mod(num1, num2 int) int {
	return num1 % num2
}

func sub(num1, num2 int) int {
	return num1 - num2
}

func ride(num1, num2 int) int {
	return num1 * num2
}

func divide(num1, num2 int) int {
	return int(math.Ceil(float64(num1) / float64(num2)))
}

func explode(s, sep string) []string {
	return strings.Split(s, sep)
}

func noescape(s string) template.HTML {
	return template.HTML(s)
}
