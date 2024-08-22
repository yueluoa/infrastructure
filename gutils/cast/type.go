package cast

import (
	"fmt"
	"reflect"
	"strconv"
)

func toInt(v reflect.Value, targetKind reflect.Kind) (i int, err error) {
	switch v.Kind() {
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
		i = int(v.Int())
	case reflect.Float64, reflect.Float32:
		i = int(v.Float())
	case reflect.String:
		if i, err = strconv.Atoi(v.String()); err != nil {
			return i, err
		}
	default:
		err = fmt.Errorf("无法将%v转成%v类型", v.Kind(), targetKind)
	}

	return
}

func toString(v reflect.Value, targetKind reflect.Kind) (s string, err error) {
	switch v.Kind() {
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
		s = strconv.Itoa(int(v.Int()))
	case reflect.Float64, reflect.Float32:
		s = fmt.Sprintf("%f", v.Float())
	case reflect.String:
		s = v.String()
	case reflect.Bool:
		s = strconv.FormatBool(v.Bool())
	default:
		err = fmt.Errorf("无法将%v转成%v类型", v.Kind(), targetKind)
	}

	return
}

func toFloat64(v reflect.Value, targetKind reflect.Kind) (f float64, err error) {
	switch v.Kind() {
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
		f = float64(v.Int())
	case reflect.Float64, reflect.Float32:
		f = v.Float()
	case reflect.String:
		if f, err = strconv.ParseFloat(v.String(), 64); err != nil {
			return
		}
	default:
		err = fmt.Errorf("无法将%v转成%v类型", v.Kind(), targetKind)
	}

	return
}

func toBool(v reflect.Value, targetKind reflect.Kind) (b bool, err error) {
	switch v.Kind() {
	case reflect.String:
		if b, err = strconv.ParseBool(v.String()); err != nil {
			return b, err
		}
	default:
		err = fmt.Errorf("无法将%v转成%v类型", v.Kind(), targetKind)
	}

	return
}

func buildIntType(targetType reflect.Type, v int) any {
	switch targetType.Kind() {
	case reflect.Int:
		return v
	case reflect.Int8:
		return int8(v)
	case reflect.Int16:
		return int16(v)
	case reflect.Int32:
		return int32(v)
	case reflect.Int64:
		return int64(v)
	case reflect.Uint:
		return uint(v)
	case reflect.Uint8:
		return uint8(v)
	case reflect.Uint16:
		return uint16(v)
	case reflect.Uint32:
		return uint32(v)
	case reflect.Uint64:
		return uint64(v)
	default:
		return v
	}
}

func buildFloatType(targetType reflect.Type, v float64) any {
	switch targetType.Kind() {
	case reflect.Float64:
		return v
	case reflect.Float32:
		return float32(v)
	default:
		return v
	}
}
