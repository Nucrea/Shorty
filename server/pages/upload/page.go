package upload

import (
	"embed"
	"text/template"

	"github.com/gin-gonic/gin"
)

//go:embed page.html
var htmlFile embed.FS

type TemplateParams struct {
	Error string
}

func NewPage() *Page {
	pageTemplate := template.Must(template.ParseFS(htmlFile, "page.html"))
	return &Page{
		pageTemplate: pageTemplate,
	}
}

type Page struct {
	pageTemplate *template.Template
}

func (p *Page) Clean(c *gin.Context) {
	c.Status(200)
	c.Header("Content-Type", "text/html")
	p.pageTemplate.Execute(c.Writer, TemplateParams{})
}

func (p *Page) WithError(c *gin.Context, error string) {
	c.Status(200)
	c.Header("Content-Type", "text/html")
	p.pageTemplate.Execute(c.Writer, TemplateParams{Error: error})
}
