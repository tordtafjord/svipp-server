package httputil

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"svipp-server/assets"
)

var Tmpl *template.Template

func init() {
	var err error
	Tmpl, err = parseTemplates()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func parseTemplates() (*template.Template, error) {
	subFS, err := fs.Sub(assets.EmbeddedFiles, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to create sub-filesystem: %w", err)
	}
	tmpl := template.Must(template.ParseFS(subFS, "*.gohtml"))

	return tmpl, nil
}
