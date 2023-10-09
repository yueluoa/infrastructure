package gutils

import (
	"math/rand"
	"reflect"
	"time"
)

var slice Slice

type Slice []interface{}

func NewSlice() Slice {
	once.Do(func() {
		slice = []interface{}{}
	})

	return slice
}

// 并集
func (s Slice) Union(vs ...interface{}) Slice {
	interfaces, hashInterfaces := getInterfaces(vs...)
	ret := make([]interface{}, 0, len(interfaces))
	for _, v := range interfaces {
		if _, ok := hashInterfaces[v]; ok {
			if hashInterfaces[v] > len(vs)-1 {
				ret = append(ret, v)
				delete(hashInterfaces, v)
				continue
			} else {
				ret = append(ret, v)
			}
		}
	}

	return ret
}

// 差集
func (s Slice) DifferenceSet(vs ...interface{}) Slice {
	interfaces, hashInterfaces := getInterfaces(vs...)
	ret := make([]interface{}, 0, len(interfaces))
	for _, v := range interfaces {
		if _, ok := hashInterfaces[v]; ok {
			if hashInterfaces[v] > len(vs)-1 {
				continue
			} else {
				ret = append(ret, v)
			}
		}
	}

	return ret
}

// 交集
func (s Slice) Intersection(vs ...interface{}) Slice {
	interfaces, hashInterfaces := getInterfaces(vs...)
	ret := make([]interface{}, 0, len(interfaces))
	for _, v := range interfaces {
		if _, ok := hashInterfaces[v]; ok {
			if hashInterfaces[v] > len(vs)-1 {
				ret = append(ret, v)
				delete(hashInterfaces, v)
			}
		}
	}

	return ret
}

// 去重
func (s Slice) RemoveDuplicate(vs interface{}) Slice {
	if reflect.TypeOf(vs).Kind() != reflect.Slice {
		return s
	}
	va := reflect.ValueOf(vs)

	ret := make([]interface{}, 0, va.Len())
	hash := map[interface{}]struct{}{}
	for i := 0; i < va.Len(); i++ {
		v := va.Index(i).Interface()
		if _, ok := hash[v]; !ok {
			hash[v] = struct{}{}
			ret = append(ret, v)
		}
	}

	return ret
}

// 随机获取一个值
func (s Slice) GetRandomValue(vs interface{}) interface{} {
	if reflect.TypeOf(vs).Kind() != reflect.Slice {
		return s
	}
	va := reflect.ValueOf(vs)

	rand.Seed(time.Now().UnixNano())

	index := rand.Intn(va.Len())
	var res interface{}
	for i := 0; i < va.Len(); i++ {
		v := va.Index(i).Interface()
		if i == index {
			res = v
		}
	}

	return res
}

// 排序 ide:= "<" || ">"
func (s Slice) Sort(ide string, vs interface{}) Slice {
	integers := make([]int64, 0)
	if reflect.TypeOf(vs).Kind() != reflect.Slice {
		return s
	}
	va := reflect.ValueOf(vs)

	chars := make([]string, 0)
	for i := 0; i < va.Len(); i++ {
		if va.Index(i).Type().Kind() == reflect.String {
			chars = append(chars, va.Index(i).String())
		} else {
			integers = append(integers, getValue(va.Index(i)))
		}
	}

	var ret []interface{}
	if len(chars) > 0 {
		ret = getCharSort(ide, chars)
	}

	if len(integers) > 0 {
		ret = getIntegerSort(ide, va.Index(0).Interface(), integers)
	}

	return ret
}

func (s Slice) String() []string {
	ret := make([]string, 0)
	for _, v := range s {
		if i, ok := v.(string); ok {
			ret = append(ret, i)
		}
	}

	return ret
}

func (s Slice) Int() []int {
	ret := make([]int, 0)
	for _, v := range s {
		if i, ok := v.(int); ok {
			ret = append(ret, i)
		}
	}

	return ret
}

func (s Slice) Int32() []int32 {
	ret := make([]int32, 0)
	for _, v := range s {
		if i, ok := v.(int32); ok {
			ret = append(ret, i)
		}
	}

	return ret
}

func (s Slice) Int64() []int64 {
	ret := make([]int64, 0)
	for _, v := range s {
		if i, ok := v.(int64); ok {
			ret = append(ret, i)
		}
	}

	return ret
}

func getIntegerSort(ide string, t interface{}, integers []int64) []interface{} {
	var ret []interface{}
	integerSort(ide, integers)
	switch t.(type) {
	case int:
		for _, v := range integers {
			ret = append(ret, int(v))
		}
	case int32:
		for _, v := range integers {
			ret = append(ret, int32(v))
		}
	default:
		for _, v := range integers {
			ret = append(ret, v)
		}
	}

	return ret
}

func integerSort(ide string, arr []int64) {
	for i := 1; i < len(arr); i++ {
		v := arr[i]
		index := i - 1
		if ide == ">" {
			for index >= 0 && arr[index] < v {
				arr[index+1] = arr[index]
				index--
			}
		}
		if ide == "<" {
			for index >= 0 && arr[index] > v {
				arr[index+1] = arr[index]
				index--
			}
		}
		if index+1 != i {
			arr[index+1] = v
		}
	}
}

func getCharSort(ide string, chars []string) []interface{} {
	var ret []interface{}
	if len(chars) == 0 {
		return ret
	}
	charSort(ide, chars)
	for _, v := range chars {
		ret = append(ret, v)
	}
	return ret
}

func charSort(ide string, arr []string) {
	for i := 1; i < len(arr); i++ {
		v := arr[i]
		index := i - 1
		if ide == ">" {
			for index >= 0 && arr[index] < v {
				arr[index+1] = arr[index]
				index--
			}
		}
		if ide == "<" {
			for index >= 0 && arr[index] > v {
				arr[index+1] = arr[index]
				index--
			}
		}
		if index+1 != i {
			arr[index+1] = v
		}
	}
}

func getValue(v reflect.Value) int64 {
	var res int64
	switch v.Interface().(type) {
	case int, int8, int16, int32, int64:
		res = v.Int()
	default:

	}

	return res
}

func getInterfaces(vs ...interface{}) ([]interface{}, map[interface{}]int) {
	interfaces := make([]interface{}, 0)
	hashInterfaces := make(map[interface{}]int)
	for _, v := range vs {
		if reflect.TypeOf(v).Kind() != reflect.Slice {
			continue
		}
		va := reflect.ValueOf(v)
		for j := 0; j < va.Len(); j++ {
			iv := va.Index(j).Interface()
			interfaces = append(interfaces, iv)
			hashInterfaces[iv]++
		}
	}

	return interfaces, hashInterfaces
}
