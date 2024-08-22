package slice

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

// Contain element是否包含target在切片中
func Contain[T comparable](items []T, element T) bool {
	for _, item := range items {
		if item == element {
			return true
		}
	}

	return false
}

// ContainFunc fn是否包含target在切片中
func ContainFunc[T any](items []T, fn func(T) bool) bool {
	for _, item := range items {
		if fn(item) {
			return true
		}
	}

	return false
}

// Keys 将切片转成map
func Keys[T any, K comparable](items []T, fn func(T) K, condFunc ...func(T) bool) map[K]T {
	result := make(map[K]T)

	if len(items) <= 0 {
		return result
	}
	for _, item := range items {
		ok := true
		for _, v := range condFunc {
			if !v(item) {
				ok = false
				break
			}
		}
		if ok {
			key := fn(item)
			result[key] = item
		}
	}

	return result
}

// GroupKeys 将切片的值进行分组
func GroupKeys[T any, K comparable](items []T, fn func(T) K, condFunc ...func(T) bool) map[K][]T {
	result := make(map[K][]T)

	if len(items) <= 0 {
		return result
	}
	for _, item := range items {
		ok := true
		for _, v := range condFunc {
			if !v(item) {
				ok = false
				break
			}
		}
		if ok {
			key := fn(item)
			if _, ok := result[key]; !ok {
				result[key] = []T{}
			}
			result[key] = append(result[key], item)
		}
	}

	return result
}

// Values 获取切片值
func Values[T, V any](items []T, fn func(T) V, condFunc ...func(T) bool) []V {
	var result []V

	if len(items) <= 0 {
		return result
	}
	for _, item := range items {
		ok := true
		for _, v := range condFunc {
			if !v(item) {
				ok = false
				break
			}
		}
		if ok {
			result = append(result, fn(item))
		}
	}

	return result
}

// Concat 多个切片合成一个新切片
func Concat[T any](slices ...[]T) []T {
	var result []T

	for _, items := range slices {
		result = append(result, items...)
	}

	return result
}

// Unique 切片元素去重
func Unique[T comparable](items []T) []T {
	var result []T

	exists := map[T]struct{}{}
	for _, item := range items {
		if _, ok := exists[item]; ok {
			continue
		}
		exists[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

// Difference 切片元素差集
func Difference[T comparable](slices ...[]T) []T {
	elementCount := make(map[T]int)

	for _, items := range slices {
		for _, item := range items {
			elementCount[item]++
		}
	}

	var result []T
	for item, count := range elementCount {
		if count == 1 {
			result = append(result, item)
		}
	}

	return result
}

// Union 切片元素并集
func Union[T comparable](slices ...[]T) []T {
	var result []T

	contain := map[T]struct{}{}
	for _, items := range slices {
		for _, item := range items {
			if _, ok := contain[item]; !ok {
				contain[item] = struct{}{}
				result = append(result, item)
			}
		}
	}

	return result
}

// Intersection 切片元素交集
func Intersection[T comparable](slices ...[]T) []T {
	var result []T

	if len(slices) <= 0 {
		return result
	}

	counts := make(map[T]int)
	for _, slice := range slices {
		exists := map[T]struct{}{}
		for _, item := range slice {
			if _, ok := exists[item]; !ok {
				counts[item]++
				exists[item] = struct{}{}
			}
		}
	}

	for item, count := range counts {
		if count == len(slices) {
			result = append(result, item)
		}
	}

	return result
}

// Equal 检查两个切片是否相等
func Equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// EqualFunc 检查两个切片是否相等
func EqualFunc[T any](a []T, b []T, fn func(T, T) bool) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if !fn(v, b[i]) {
			return false
		}
	}

	return true
}

// Filter 过滤有效值组成新切片
func Filter[T any](items []T, fn func(int, T) bool) []T {
	result := make([]T, 0)

	for i, v := range items {
		if fn(i, v) {
			result = append(result, v)
		}
	}

	return result
}

// Count 统计切片中指定值出现次数
func Count[T comparable](items []T, item T) int {
	var count int

	for _, v := range items {
		if item == v {
			count++
		}
	}

	return count
}

// CountFunc 统计切片中指定值出现次数
func CountFunc[T any](items []T, fn func(int, T) bool) int {
	var count int

	for i, v := range items {
		if fn(i, v) {
			count++
		}
	}

	return count
}

// IndexOf 返回切片中指定值索引
func IndexOf[T comparable](items []T, element T) int {
	for i, e := range items {
		if e == element {
			return i
		}
	}

	return -1
}

// LastIndexOf 返回切片中指定值最后索引
func LastIndexOf[T comparable](items []T, element T) int {
	for i := len(items) - 1; i >= 0; i-- {
		if element == items[i] {
			return i
		}
	}

	return -1
}

// ToSlicePointer 返回切片的指针
func ToSlicePointer[T any](items ...T) []*T {
	result := make([]*T, len(items))

	for i := range items {
		result[i] = &items[i]
	}

	return result
}

// ToSlice 转成切片
func ToSlice[T any](items ...T) []T {
	result := make([]T, len(items))
	copy(result, items)

	return result
}

// ToMaps 将结构体切片转换为切片，每个元素是一个 map[string]any
func ToMaps[T any](items []T) []map[string]any {
	var result []map[string]any

	for _, item := range items {
		v := reflect.ValueOf(item)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		t := v.Type()
		m := make(map[string]any)
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			if field.PkgPath != "" {
				continue
			}
			tag := field.Tag.Get("json")
			if tag == "" {
				tag = field.Name
			}
			m[tag] = value.Interface()
		}

		result = append(result, m)
	}

	return result
}

// Append 追加元素
func Append[T comparable](items []T, element T) []T {
	items = append(items, element)
	return items
}

// AppendIfAbsent 不存在则追加元素
func AppendIfAbsent[T comparable](items []T, element T) []T {
	if !Contain(items, element) {
		items = append(items, element)
	}

	return items
}

// Join 将切片转换为字符串，其中元素由指定的分隔符分隔
func Join[T any](items []T, delimiter string) string {
	var sb strings.Builder

	for i, item := range items {
		if i > 0 {
			sb.WriteString(delimiter)
		}
		sb.WriteString(fmt.Sprint(item))
	}

	return sb.String()
}

// ForEach 将给定函数应用于切片的每个元素
func ForEach[T any](items []T, fn func(T)) {
	for _, item := range items {
		fn(item)
	}
}

// Replace 切片新值替换部分旧值
func Replace[T comparable](items []T, old T, new T, n int) []T {
	result := make([]T, len(items))
	copy(result, items)

	for i := range result {
		if result[i] == old && n != 0 {
			result[i] = new
			n--
		}
	}

	return result
}

// ReplaceAll 切片新值替换全部旧值
func ReplaceAll[T comparable](items []T, old T, new T) []T {
	return Replace(items, old, new, -1)
}

// DeleteAt 删除切片元素
func DeleteAt[T any](items []T, index int) []T {
	if index >= len(items) {
		index = len(items) - 1
	}
	if len(items) <= 0 {
		return items
	}

	result := make([]T, len(items)-1)
	copy(result, items[:index])
	copy(result[index:], items[index+1:])

	return result
}

// Paginate 切片分页
func Paginate[T comparable](items []T, page, pageSize int) []T {
	if page <= 0 {
		page = 1
	}

	size := len(items)
	if pageSize <= 0 || size <= 0 {
		return items
	}
	offset := (page - 1) * pageSize
	limit := offset + pageSize

	if limit > size {
		limit = size
	}
	if offset >= size {
		offset = size % pageSize
	}

	return items[offset:limit]
}

// UpdateAt 更新索引处的切片元素
func UpdateAt[T any](items []T, index int, value T) []T {
	size := len(items)

	if index < 0 || index >= size {
		return items
	}
	items = append(items[:index], append([]T{value}, items[index+1:]...)...)

	return items
}

// InsertAt 将值或其他切片插入索引处的切片中
func InsertAt[T any](items []T, index int, value any) []T {
	size := len(items)

	if index < 0 || index > size {
		return items
	}

	if v, ok := value.(T); ok {
		items = append(items[:index], append([]T{v}, items[index:]...)...)
		return items
	}

	if v, ok := value.([]T); ok {
		items = append(items[:index], append(v, items[index:]...)...)
		return items
	}

	return items
}

// IsAscending 切片是否按照升序排序
func IsAscending[T constraints.Ordered](items []T) bool {
	for i := 1; i < len(items); i++ {
		if items[i-1] > items[i] {
			return false
		}
	}

	return true
}

// IsDescending 切片是否按照降序排序
func IsDescending[T constraints.Ordered](items []T) bool {
	for i := 1; i < len(items); i++ {
		if items[i-1] < items[i] {
			return false
		}
	}

	return true
}

// IsSorted 切片是否排序
func IsSorted[T constraints.Ordered](items []T) bool {
	return IsAscending(items) || IsDescending(items)
}

// Sort 切片排序
func Sort[T any](items []T, less func(T, T) bool) {
	if len(items) <= 0 {
		return
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})
}

// Random 获取切片随机元素
func Random[T any](items []T) T {
	var result T

	if len(items) <= 0 {
		return result
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result = items[r.Intn(len(items))]

	return result
}

// Transform 切片转换字段类型
func Transform[T any, U any](items []T, fn func(T) U) []U {
	result := make([]U, len(items))

	for i, v := range items {
		result[i] = fn(v)
	}

	return result
}

// Sum 切片值相加
func Sum[T any, V constraints.Integer](items []T, fn func(T) V) V {
	var result V

	for _, v := range items {
		result += fn(v)
	}

	return result
}
