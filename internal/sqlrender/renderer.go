package sqlrender

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed templates/*.sql.tmpl
var templatesFS embed.FS

type Renderer struct {
	t *template.Template
}

func NewRenderer() (*Renderer, error) {
	funcs := template.FuncMap{
		"qident":           qident,
		"fqSchema":         fqSchema,
		"fqTable":          fqTable,
		"fqTableFromTable": fqTableFromTable,
		"s3Join":           s3Join,
		"sqlString":        sqlString,
		"sqlValue":         sqlValue,
		"propsSQL":         propsSQL,
		"tableWithProps":   tableWithProps,
		"notLastColumn":    notLastColumn,
	}

	t, err := template.New("").Funcs(funcs).ParseFS(templatesFS, "templates/*.sql.tmpl")
	if err != nil {
		return nil, err
	}

	return &Renderer{t: t}, nil
}

func (r *Renderer) Render(name string, data any) (string, error) {
	var buf bytes.Buffer

	err := r.t.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
