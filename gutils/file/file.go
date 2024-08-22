package file

import (
	"io"
	"mime"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// CreateFile 创建文件(目录需要存在&文件存在也会重新创建变成空文件)
func CreateFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	return nil
}

// CreateFileIfNotExists 如果文件不存在则创建一个新文件
func CreateFileIfNotExists(path string) error {
	dir := filepath.Dir(path)
	// 检查目录是否存在，不存在则创建
	if err := CreateDir(dir); err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()
	} else if err != nil {
		return err
	}

	return nil
}

// CopyFile 将源文件复制到目标文件
func CopyFile(srcPath string, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	dist, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() { _ = dist.Close() }()

	_, err = io.Copy(dist, src)

	return err
}

type File struct {
	Name      string
	Size      int64
	Mode      os.FileMode
	ModTime   time.Time
	IsDir     bool
	MediaType string
}

// GetFileInfo 文件信息
func GetFileInfo(path string) (*File, error) {
	var f = &File{}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return f, err
	}
	ext := filepath.Ext(path)
	mediaType := mime.TypeByExtension(ext)

	f.Name = fileInfo.Name()
	f.Size = fileInfo.Size()
	f.Mode = fileInfo.Mode()
	f.ModTime = fileInfo.ModTime()
	f.IsDir = fileInfo.IsDir()
	f.MediaType = mediaType

	return f, nil
}

// RemoveFile 删除文件
func RemoveFile(path string) error {
	return os.Remove(path)
}

// ClearFile 清空文件
func ClearFile(path string) error {
	return os.Truncate(path, 0)
}

// CurrentPath 返回当前绝对路径
func CurrentPath() string {
	var absPath string

	_, filename, _, ok := runtime.Caller(1)
	if ok {
		absPath = filepath.Dir(filename)
	}

	return absPath
}

// CreateDir 创建文件夹(如果path已经是一个目录不执行任何操作)
func CreateDir(absPath string) error {
	return os.MkdirAll(absPath, os.ModePerm)
}

// RemoveDir 移除文件夹
func RemoveDir(path string) error {
	return os.RemoveAll(path)
}
