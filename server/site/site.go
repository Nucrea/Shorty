package site

import (
	"shorty/server/site/pages"

	"github.com/gin-gonic/gin"
)

type Site struct{}

func (s *Site) err(c *gin.Context, status int, msg string) {
	pages.ErrorPage(status, msg).Render(c, c.Writer)
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
	pages.LinkForm().Render(c, c.Writer)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) ImageForm(c *gin.Context, id, captchaBase64 string) {
	pages.ImageForm(id, captchaBase64).Render(c, c.Writer)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) FileForm(c *gin.Context, id, captchaBase64 string) {
	pages.FileForm(id, captchaBase64).Render(c, c.Writer)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) LinkResult(c *gin.Context, url string, qrBase64 string) {
	pages.LinkResult(pages.LinkResultParams{Shortlink: url, QRBase64: qrBase64}).Render(c, c.Writer)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) ImageView(c *gin.Context, p pages.ImageViewParams) {
	pages.ImageView(p).Render(c, c.Writer)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) FileView(c *gin.Context, p pages.FileViewParams) {
	pages.FileView(p).Render(c, c.Writer)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}

func (s *Site) FileDownload(c *gin.Context, fileRawUrl string) {
	pages.FileDownload(fileRawUrl).Render(c, c.Writer)
	c.Header("Content-Type", "text/html")
	c.Status(200)
}
