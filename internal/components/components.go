package components

import (
	"embed"
	"io/fs"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"text/template"
)

const (
	templatesDir = "templateAssets"
)

var (
	//go:embed templateAssets/*.yaml
	files     embed.FS
	templates map[string]*template.Template
)

func LoadTemplates() error {
	// Templates is a "singleton"
	if templates != nil {
		return nil
	}

	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name())
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}

func Get(filename string, unstructured bool, data interface{}) (client.Object, error) {
	err := LoadTemplates()
	if err != nil {
		return nil, err
	}

	out, err := renderTemplate(templates[filename], data)
	if err != nil {
		return nil, err
	}
	var obj client.Object
	if unstructured {
		obj, err = parseUnstructured(out)
	} else {
		obj, err = parseObject(out)
	}
	if err != nil {
		return nil, err
	}
	return obj, nil
}
