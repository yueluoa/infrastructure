package has

import (
	"fmt"
	"reflect"
)

// Keys 返回has键切片
func Keys[K comparable, V any](m map[K]V, condFunc ...func(K) bool) []K {
	keys := make([]K, 0, len(m))

	if len(m) <= 0 {
		return keys
	}

	for k := range m {
		ok := true
		for _, v := range condFunc {
			if !v(k) {
				ok = false
				break
			}
		}
		if ok {
			keys = append(keys, k)
		}
	}

	return keys
}

// KeysFunc 返回has键切片
func KeysFunc[K comparable, V, T any](m map[K]V, fn func(K) T, condFunc ...func(K) bool) []T {
	keys := make([]T, 0, len(m))

	if len(m) <= 0 {
		return keys
	}

	for k := range m {
		ok := true
		for _, v := range condFunc {
			if !v(k) {
				ok = false
				break
			}
		}
		if ok {
			keys = append(keys, fn(k))
		}
	}

	return keys
}

// Value 返回has值
func Value[K comparable, T any](m map[K]T, key K, defaultValue T) T {
	value, ok := m[key]
	if !ok {
		return defaultValue
	}

	return value
}

// Values 返回has值切片
func Values[K comparable, V any](m map[K]V, condFunc ...func(V) bool) []V {
	values := make([]V, 0, len(m))

	if len(m) <= 0 {
		return values
	}

	for _, item := range m {
		ok := true
		for _, v := range condFunc {
			if !v(item) {
				ok = false
				break
			}
		}
		if ok {
			values = append(values, item)
		}
	}

	return values
}

// ValuesFunc 返回has值切片
func ValuesFunc[K comparable, V, T any](m map[K]V, fn func(V) T, condFunc ...func(V) bool) []T {
	values := make([]T, 0, len(m))

	if len(m) <= 0 {
		return values
	}

	for _, item := range m {
		ok := true
		for _, v := range condFunc {
			if !v(item) {
				ok = false
				break
			}
		}
		if ok {
			values = append(values, fn(item))
		}
	}

	return values
}

// Merge 合并(key相同则后面覆盖前面)
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}

// ForEach 将给定函数应用于has的每个元素
func ForEach[K comparable, V any](m map[K]V, fn func(K, V)) {
	for k, v := range m {
		fn(k, v)
	}
}

// Filter 过滤map
func Filter[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	result := make(map[K]V)

	for k, v := range m {
		if fn(k, v) {
			result[k] = v
		}
	}

	return result
}

// ToStruct 映射到结构体
func ToStruct[T any](m map[string]any, structObj *T) error {
	v := reflect.ValueOf(structObj).Elem()

	for key, value := range m {
		field := v.FieldByName(key)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		val := reflect.ValueOf(value)
		if field.Type() != val.Type() {
			if val.Type().ConvertibleTo(field.Type()) {
				val = val.Convert(field.Type())
			} else {
				return fmt.Errorf("type mismatch for field %s", key)
			}
		}
		field.Set(val)
	}

	return nil
}
