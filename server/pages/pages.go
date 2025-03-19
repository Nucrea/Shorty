package pages

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

func (s *Site) LinkForm(c *gin.Context) {
	s.template("views/link_form.html").Execute(c.Writer, nil)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) ImageForm(c *gin.Context, id, captchabase64 string) {
	s.template("views/image_form.html").Execute(c.Writer, ImageFormParams{Id: id, CaptchaBase64: captchabase64})
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) FileForm(c *gin.Context, id, captchabase64 string) {
	s.template("views/file_form.html").Execute(c.Writer, LinkFormParams{Id: id, CaptchaBase64: captchabase64})
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) LinkResult(c *gin.Context, url string, qrBase64 string) {
	s.template("views/link_result.html").Execute(c.Writer, LinkResultParams{Shortlink: url, QRBase64: qrBase64})
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) ImageView(c *gin.Context, p ViewImageParams) {
	s.template("views/image_view.html").Execute(c.Writer, p)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) FileView(c *gin.Context, p ViewFileParams) {
	s.template("views/file_view.html").Execute(c.Writer, p)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}
