package webutil

import (
	"fmt"
	"github.com/hoisie/web"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type TemplateSet struct {
	root *template.Template
}

func (ts *TemplateSet) AddDirectory(root string) []error {
	allErrors := []error{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if info.IsDir() {
				return filepath.SkipDir
			}
			allErrors = append(allErrors, fmt.Errorf("%s: %s", path, err))
			return nil
		}
		if info.IsDir() {
			return nil
		}
		defer func() {
			if err != nil {
				allErrors = append(allErrors, fmt.Errorf("%s: %s", path, err))
			}
		}()
		name, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		name = filepath.ToSlash(name)
		fmt.Println(root, path, name)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		if ts.root == nil {
			ts.root, err = template.New(name).Parse(string(content))
		} else {
			_, err = ts.root.New(name).Parse(string(content))
		}
		return nil
	})
	if len(allErrors) == 0 {
		return nil
	}
	return allErrors
}

func (ts *TemplateSet) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return ts.root.ExecuteTemplate(wr, name, data)
}

func (ts *TemplateSet) ExecuteTemplateWithContext(ctx *web.Context, name string, data interface{}) {
	if ts.root.Lookup(name) == nil {
		ctx.NotFound("Template " + name + " not found")
		return
	}
	err := ts.ExecuteTemplate(ctx, name, data)
	if err != nil {
		Error(ctx, err)
		return
	}
	ext := filepath.Ext(name)
	if ext == "" {
		ext = ".html"
	}
	ctx.ContentType(ext)
}

func (ts *TemplateSet) HandleTemplate(name string, data interface{}) func(*web.Context) {
	return func(ctx *web.Context) {
		ts.ExecuteTemplateWithContext(ctx, name, data)
	}
}

var defaultTemplateSet = new(TemplateSet)

func init() {
	defaultTemplateSet.AddDirectory("templates")
}

func ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return defaultTemplateSet.ExecuteTemplate(wr, name, data)
}

func ExecuteTemplateWithContext(ctx *web.Context, name string, data interface{}) {
	defaultTemplateSet.ExecuteTemplateWithContext(ctx, name, data)
}

func HandleTemplate(name string, data interface{}) func(*web.Context) {
	return defaultTemplateSet.HandleTemplate(name, data)
}
