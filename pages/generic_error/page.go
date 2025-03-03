package genericerror

import (
	"embed"
	"text/template"

	"github.com/gin-gonic/gin"
)

//go:embed page.html
var htmlFile embed.FS

type TemplateParams struct {
	Title string
	Text  string
}

func NewPage() *Page {
	pageTemplate := template.Must(template.ParseFS(htmlFile, "index.html"))
	return &Page{
		pageTemplate: pageTemplate,
	}
}

type Page struct {
	pageTemplate *template.Template
}

func (p *Page) NotFound(c *gin.Context) {
	p.pageTemplate.Execute(c.Writer, TemplateParams{
		Title: "Not found",
		Text:  "404 Not Found",
	})
	c.Status(404)
	c.Header("Content-Type", "text/html")
}

func (p *Page) InternalError(c *gin.Context) {
	p.pageTemplate.Execute(c.Writer, TemplateParams{
		Title: "Internal Error",
		Text:  "500 Internal Error",
	})
	c.Status(500)
	c.Header("Content-Type", "text/html")
}
