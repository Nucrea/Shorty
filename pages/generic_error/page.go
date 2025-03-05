package genericerror

import (
	"embed"
	"text/template"

	"github.com/gin-gonic/gin"
)

//go:embed page.html
var htmlFile embed.FS

type TemplateParams struct {
	Status int
	Text   string
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

func (p *Page) NotFound(c *gin.Context) {
	p.pageTemplate.Execute(c.Writer, TemplateParams{
		Status: 404,
		Text:   "Not Found",
	})
	c.Header("Content-Type", "text/html")
	c.Status(500)
}

func (p *Page) InternalError(c *gin.Context) {
	p.pageTemplate.Execute(c.Writer, TemplateParams{
		Status: 404,
		Text:   "Internal Error",
	})
	c.Header("Content-Type", "text/html")
	c.Status(500)
}
