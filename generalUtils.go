package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/masterminds/sprig"
)

func ParseAllTemplates(location string) (*template.Template, error) {
	templatesVar := template.New("").Funcs(sprig.FuncMap())
	filepath.WalkDir(location, func(path string, _ fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".html") {
			_, err := templatesVar.New("").Funcs(sprig.FuncMap()).ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return templatesVar, nil
}
