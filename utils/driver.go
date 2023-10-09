package utils

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/yueluoa/infrastructure/gerror"
)

type SliceUint8 []uint8

func (s *SliceUint8) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceUint8) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type SliceUint16 []uint16

func (s *SliceUint16) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceUint16) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type SliceUint32 []uint32

func (s *SliceUint32) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceUint32) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type SliceUint64 []uint64

func (s *SliceUint64) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceUint64) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type SliceUint []uint

func (s *SliceUint) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceUint) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type SliceInt8 []int8

func (s *SliceInt8) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceInt8) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type SliceInt16 []int16

func (s *SliceInt16) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceInt16) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type SliceInt32 []int32

func (s *SliceInt32) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceInt32) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type SliceInt64 []int64

func (s *SliceInt64) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceInt64) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type SliceInt []int

func (s *SliceInt) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceInt) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type SliceString []string

func (s *SliceString) Scan(val interface{}) error {
	b, _ := val.([]byte)
	return json.Unmarshal(b, &s)
}

func (s SliceString) Value() (driver.Value, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

type CustomBinary []byte

func (cb *CustomBinary) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return gerror.New("Failed to scan CustomBinary value")
	}
	*cb = append([]byte{}, bytes...)
	return nil
}

func (cb CustomBinary) Value() (driver.Value, error) {
	return cb, nil
}
