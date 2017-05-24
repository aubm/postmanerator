package themes

import (
	"io"
	"text/template"

	"github.com/aubm/postmanerator/postman"
)

const (
	templateName  = ""
	mainThemeFile = "index.tpl"
)

type Renderer struct{}

func (r *Renderer) Render(w io.Writer, theme *Theme, collection postman.Collection) error {
	tmpl, err := template.New(templateName).Funcs(r.getTemplateHelpers()).ParseFiles(theme.Files...)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, mainThemeFile, collection)
}

func (r *Renderer) getTemplateHelpers() template.FuncMap {
	return template.FuncMap{
		"curlSnippet":  curlSnippet,
		"findResponse": helperFindResponse,
		"hasContent":   helperHasContent,
		"httpSnippet":  helperHttpSnippet,
		"indentJSON":   helperIndentJSON,
		"inline":       helperInline,
		"markdown":     helperMarkdown,
		"slugify":      helperSlugify,
	}
}
