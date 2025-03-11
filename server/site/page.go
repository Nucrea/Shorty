package site

import (
	"embed"
	"html/template"

	"github.com/gin-gonic/gin"
)

//go:embed layout.html views
var views embed.FS

type Site struct{}

func (s *Site) template(view string) *template.Template {
	tmpl, err := template.ParseFS(views, "layout.html", view)
	if err != nil {
		panic(err)
	}
	return tmpl
}

func (s *Site) err(c *gin.Context, status int, msg string) {
	s.template("views/error.html").Execute(c.Writer, ErrParams{Status: status, Message: msg})
	c.Header("Content-Type", "text/html")
	c.AbortWithStatus(status)
}

func (s *Site) TemporarilyBanned(c *gin.Context) {
	s.err(c, 403, "Temporarily Banned")
}

func (s *Site) TooManyRequests(c *gin.Context) {
	s.err(c, 429, "Too Many Requests")
}

func (s *Site) NotFound(c *gin.Context) {
	s.err(c, 404, "Not Found")
}

func (s *Site) InternalError(c *gin.Context) {
	s.err(c, 500, "Internal Error")
}

func (s *Site) CreateLink(c *gin.Context) {
	s.template("views/create_link.html").Execute(c.Writer, nil)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) UploadImage(c *gin.Context) {
	s.template("views/upload_image.html").Execute(c.Writer, nil)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) LinkResult(c *gin.Context, url string) {
	s.template("views/link_result.html").Execute(c.Writer, LinkResultParams{Url: url})
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) QRResult(c *gin.Context, imageBase64 string) {
	s.template("views/qr_result.html").Execute(c.Writer, QRResultParams{ImageBase64: imageBase64})
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) ViewImage(c *gin.Context, p ViewImageParams) {
	s.template("views/view_image.html").Execute(c.Writer, p)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}
