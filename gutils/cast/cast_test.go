package cast

import (
	"fmt"
	"testing"
)

func TestCast(t *testing.T) {
	result1, err := ToAny[int8]("1")
	fmt.Println(result1, err)
	result2, err := ToAny[int16]("1")
	fmt.Println(result2, err)
	result3, err := ToAny[int32]("1")
	fmt.Println(result3, err)
	result4, err := ToAny[int64]("177")
	fmt.Println(result4, err)
	result5, err := ToAny[string](199)
	fmt.Println(result5, err)
	result6, err := ToAny[bool]("true")
	fmt.Println(result6, err)
}
