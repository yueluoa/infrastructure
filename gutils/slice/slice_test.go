package slice

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSlice(t *testing.T) {
	type ab struct {
		A string
		B int
	}
	array1 := []ab{{A: "1", B: 1}, {A: "2", B: 2}, {A: "2", B: 2}}
	result1 := ContainFunc(array1, func(f ab) bool { return f.A == "1" && f.B == 1 })
	result2 := Contain(array1, ab{A: "1", B: 1})
	fmt.Println(result1)
	fmt.Println(result2)

	slice1 := []string{"a", "b", "c"}
	slice2 := []string{"d", "e", "f"}
	slice3 := []string{"g", "h", "i"}
	result3 := Concat(slice1, slice2, slice3)
	fmt.Println(result3)

	strA := []string{"apple", "banana", "cherry"}
	strB := []string{"banana", "date", "fig"}
	strC := []string{"apple", "fig", "grape"}
	result4 := Difference(strA, strB, strC)
	fmt.Println(result4)

	slice4 := []int{1, 2, 3}
	slice5 := []int{2, 4, 6}
	isDouble := func(a, b int) bool {
		return b == a*2
	}

	result5 := EqualFunc(slice4, slice5, isDouble)
	fmt.Println(result5)

	nums := []int{1, 2, 3, 4, 5}
	isEven := func(i, num int) bool {
		return num%2 == 0
	}
	result6 := Filter(nums, isEven)
	fmt.Println(result6)

	result7 := Count(nums, 2)
	fmt.Println(result7)

	result8 := CountFunc(nums, isEven)
	fmt.Println(result8)

	result9 := IndexOf(array1, ab{A: "1", B: 1})
	fmt.Println(result9)

	result10 := LastIndexOf(array1, ab{A: "2", B: 2})
	fmt.Println(result10)

	result11 := Join(nums, "")
	fmt.Println(result11)

	result12 := Replace([]string{"a", "b", "c", "a"}, "a", "x", 0)
	result13 := Replace([]string{"a", "b", "c", "a"}, "a", "x", 1)
	fmt.Println(result12)
	fmt.Println(result13)

	result14 := DeleteAt([]int{}, 1)
	fmt.Println(result14)

	result15 := Unique(nums)
	fmt.Println(result15)

	result16 := Union([]int{1, 3, 4, 6}, []int{1, 2, 5, 6})
	fmt.Println(result16)

	result17 := Random(nums)
	fmt.Println(result17)

	result18 := Keys([]string{"a", "ab", "abc"}, func(str string) int {
		return len(str)
	}, func(s string) bool {
		return s == "a"
	})
	fmt.Println(result18)

	result19 := Values([]string{"a", "ab", "abc"}, func(str string) string {
		return str
	})
	fmt.Println(result19)

	str1 := "a"
	str2 := "b"
	result20 := ToSlicePointer(str1, str2)
	fmt.Println(reflect.DeepEqual(result20, []*string{&str1, &str2}))

	result21 := ToSlice("a", "b", "c")
	fmt.Println(result21)

	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	people := []*person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	peopleMaps := ToMaps(people)
	for _, v := range peopleMaps {
		fmt.Println(v)
	}

	// 计算多个切片的交集
	result22 := Intersection([]int{1, 2, 3, 4, 5}, []int{3, 4, 5, 6, 7}, []int{5, 6, 7, 8, 9})
	fmt.Println(result22)
}
