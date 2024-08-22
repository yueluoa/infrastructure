package has

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"
)

func TestHas(t *testing.T) {
	m := map[int]string{
		1: "a",
		2: "a",
		3: "b",
		4: "c",
		5: "d",
	}

	result1 := Keys(m)
	sort.Ints(result1)
	fmt.Println(result1)

	result2 := KeysFunc(m,
		func(item int) int { return item },
		func(i int) bool { return i >= 5 })
	sort.Ints(result2)
	fmt.Println(result2)

	value := Value(m, 1, "ww")
	fmt.Println(value)

	result3 := Values(m)
	sort.Strings(result3)
	fmt.Println(result3)

	result4 := ValuesFunc(m,
		func(v string) string { return v },
		func(s string) bool { return s == "a" })
	sort.Strings(result4)
	fmt.Println(result4)

	m1 := map[int]string{
		1: "a",
		2: "b",
	}
	m2 := map[int]string{
		1: "c",
		3: "d",
	}
	result5 := Merge(m1, m2)
	fmt.Println(result5)

	var sum string
	ForEach(m, func(i int, s string) {
		sum += s
	})
	fmt.Println(sum)

	result6 := Filter(m, func(i int, s string) bool {
		return i > 30
	})
	fmt.Println(result6)

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var p Person
	err := ToStruct(map[string]interface{}{
		"00":  "Alice",
		"Age": 30,
	}, &p)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		b, _ := json.Marshal(p)
		fmt.Println("Person struct:", string(b))
	}
}
