package html

import (
	"fmt"
	"github.com/gin-gonic/gin/render"
	"html/template"
	"path/filepath"
)

type Render map[string]*template.Template

var (
	_ render.HTMLRender = Render{}
	_ Renderer          = Render{}
)

type Renderer interface {
	render.HTMLRender
	Add(name string, tmpl *template.Template)
	AddFromFiles(name string, files ...string) *template.Template
	AddFromGlob(name, glob string) *template.Template
	AddFromString(name, templateString string) *template.Template
	AddFromStringsFuncs(name string, funcMap template.FuncMap, templateStrings ...string) *template.Template
	AddFromFilesFuncs(name string, funcMap template.FuncMap, files ...string) *template.Template
}

func New() Render {
	return make(Render)
}

func (r Render) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	if _, ok := r[name]; ok {
		panic(fmt.Sprintf("template %s already exists", name))
	}

	r[name] = tmpl
}

func (r Render) AddFromFiles(name string, files ...string) *template.Template {
	tmpl := template.Must(template.ParseFiles(files...))
	r.Add(name, tmpl)

	return tmpl
}

func (r Render) AddFromGlob(name, glob string) *template.Template {
	tmpl := template.Must(template.ParseGlob(glob))
	r.Add(name, tmpl)

	return tmpl
}

func (r Render) AddFromString(name, templateString string) *template.Template {
	tmpl := template.Must(template.New(name).Parse(templateString))
	r.Add(name, tmpl)

	return tmpl
}

func (r Render) AddFromStringsFuncs(name string, funcMap template.FuncMap, templateStrings ...string) *template.Template {
	tmpl := template.New(name).Funcs(funcMap)

	for _, ts := range templateStrings {
		tmpl = template.Must(tmpl.Parse(ts))
	}
	r.Add(name, tmpl)

	return tmpl
}

func (r Render) AddFromFilesFuncs(name string, funcMap template.FuncMap, files ...string) *template.Template {
	baseName := filepath.Base(files[0])
	tmpl := template.Must(template.New(baseName).Funcs(funcMap).ParseFiles(files...))
	r.Add(name, tmpl)

	return tmpl
}

func (r Render) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r[name],
		Data:     data,
	}
}
