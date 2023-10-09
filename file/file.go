package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// 创建文件夹
func CreateDir(path string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
	}

	return err
}

// 创建文件/文件夹不存在则先创建文件夹
func CreateFile(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

// 创建文件
func WriteFile(filePath, content string) error {
	var (
		file    *os.File
		fileErr error
	)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, fileErr = os.Create(filePath)
	} else {
		file, fileErr = os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC, 0666)
	}
	if fileErr != nil {
		return fmt.Errorf("文件操作失败: %v", fileErr)
	}
	defer file.Close()

	_, fileErr = file.WriteString(content)
	if fileErr != nil {
		return fmt.Errorf("文件写入失败: %v", fileErr)
	}

	return nil
}

// 下载文件/图片
func DownloadFile(url string, fp string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = os.MkdirAll(filepath.Dir(fp), os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	return err
}
