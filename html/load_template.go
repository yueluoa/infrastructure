package html

import (
	"os"
	"path/filepath"
	"strings"
)

func LoadTemplates(templatesDir string) Renderer {
	r := New()

	layouts, err := filepath.Glob(templatesDir + "/common/layout.html")
	if err != nil {
		panic(err)
	}
	includes, err := filepath.Glob(templatesDir + "/*.html") // 一级嵌套文件
	if err != nil {
		panic(err)
	}
	nextIncludes, err := filepath.Glob(templatesDir + "/*/*.html") // 二级嵌套文件
	if err != nil {
		panic(err)
	}
	includes = append(includes, nextIncludes...)

	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		//r.AddFromFiles(filepath.Base(include), files...)
		r.AddFromFilesFuncs(filepath.Base(include), hashFunc, files...)
	}

	return r
}

func GetTemplates(root string) []string {
	var files []string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return err
	})
	if len(files) == 0 {
		files = make([]string, 0)
	}

	return files
}
