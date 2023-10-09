package utils

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

const CommonPageSize = 10

// 对分页做兜底限制
func Paginate(db *gorm.DB, page, pageSize int) *gorm.DB {
	switch {
	case pageSize > 100 || pageSize <= 0:
		pageSize = 100
	case page <= 0:
		page = 1
	}
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize)

	return db
}

func StructToMap(in interface{}, tagName ...string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return out, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	tn := "json"
	if len(tagName) > 0 {
		tn = tagName[0]
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tn); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}

	return out, nil
}
