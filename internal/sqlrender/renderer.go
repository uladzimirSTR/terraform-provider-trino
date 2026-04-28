package sqlrender

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed templates/*.sql.tmpl
var templatesFS embed.FS

type Renderer[Model any] struct {
	t *template.Template
}

func NewRenderer[Model any]() (*Renderer[Model], error) {
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

	return &Renderer[Model]{t: t}, nil
}

func (r *Renderer[Model]) Render(name string, data Model) (string, error) {
	var buf bytes.Buffer

	err := r.t.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
