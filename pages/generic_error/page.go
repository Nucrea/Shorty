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

func (p *Page) TooMuchRequests(c *gin.Context) {
	c.Status(429)
	p.pageTemplate.Execute(c.Writer, TemplateParams{
		Status: 429,
		Text:   "Too Many Requests. Your IP is temporarily banned.",
	})
	c.Header("Content-Type", "text/html")
	c.Abort()
}

func (p *Page) NotFound(c *gin.Context) {
	c.Status(404)
	p.pageTemplate.Execute(c.Writer, TemplateParams{
		Status: 404,
		Text:   "Not Found",
	})
	c.Header("Content-Type", "text/html")
	c.Abort()
}

func (p *Page) InternalError(c *gin.Context) {
	c.Status(500)
	p.pageTemplate.Execute(c.Writer, TemplateParams{
		Status: 500,
		Text:   "Internal Error",
	})
	c.Header("Content-Type", "text/html")
	c.Abort()
}
