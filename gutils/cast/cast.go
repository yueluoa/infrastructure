package cast

import (
	"fmt"
	"reflect"
)

func ToString(value any) string {
	s, _ := ToAny[string](value)
	return s
}

func ToInt(value any) int {
	i, _ := ToAny[int](value)
	return i
}

// ToAny 适用于数字、浮点数和字符串相互之间的类型转换
func ToAny[T any](value any) (a T, err error) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	targetType := reflect.TypeOf(a)
	kind := targetType.Kind()
	switch kind {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		var i int
		if i, err = toInt(v, kind); err != nil {
			return
		}
		a = buildIntType(targetType, i).(T)
	case reflect.String:
		var s string
		if s, err = toString(v, kind); err != nil {
			return
		}
		a = any(s).(T)
	case reflect.Float64, reflect.Float32:
		var f float64
		if f, err = toFloat64(v, kind); err != nil {
			return
		}
		a = buildFloatType(targetType, f).(T)
	case reflect.Bool:
		var b bool
		if b, err = toBool(v, kind); err != nil {
			return
		}
		a = any(b).(T)
	default:
		err = fmt.Errorf("unsupported target type")
	}

	return
}
