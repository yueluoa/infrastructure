package file

import (
	"fmt"
	"testing"
)

func TestFile(t *testing.T) {
	filePath := "./a/b/aaa.json"

	//err := CreateFile(filePath)
	//fmt.Println("failed to create file: ", err)

	err := ClearFile(filePath)
	fmt.Println("failed to clear file: ", err)

	absPath := CurrentPath()
	fmt.Println(absPath)
}
