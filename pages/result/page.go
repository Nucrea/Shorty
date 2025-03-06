package result

import (
	"embed"
	"text/template"

	"github.com/gin-gonic/gin"
)

//go:embed page.html
var htmlFile embed.FS

type TemplateParams struct {
	Shortlink string
	QRCode    string
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

func (p *Page) WithQRCode(c *gin.Context, qrCodeBase64 string) {
	c.Status(200)
	c.Header("Content-Type", "text/html")
	p.pageTemplate.Execute(c.Writer, TemplateParams{QRCode: qrCodeBase64})
}

func (p *Page) WithLink(c *gin.Context, link string) {
	c.Status(200)
	c.Header("Content-Type", "text/html")
	p.pageTemplate.Execute(c.Writer, TemplateParams{Shortlink: link})
}
