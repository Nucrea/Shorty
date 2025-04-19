package site

import (
	"shorty/server/site/pages"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

type Site struct{}

func (s *Site) err(c *gin.Context, status int, msg string) templ.Component {
	c.Header("Content-Type", "text/html")
	c.AbortWithStatus(status)
	return pages.ErrorPage(status, msg)
}

func (s *Site) TemporarilyBanned(c *gin.Context) templ.Component {
	return s.err(c, 403, "Temporarily Banned")
}

func (s *Site) TooManyRequests(c *gin.Context) templ.Component {
	return s.err(c, 429, "Too Many Requests")
}

func (s *Site) NotFound(c *gin.Context) templ.Component {
	return s.err(c, 404, "Not Found")
}

func (s *Site) InternalError(c *gin.Context) templ.Component {
	return s.err(c, 500, "Internal Error")
}

func (s *Site) LinkForm(c *gin.Context) templ.Component {
	c.Header("Content-Type", "text/html")
	c.Status(200)
	return pages.LinkForm()
}

func (s *Site) ImageForm(c *gin.Context, id, captchaBase64 string) templ.Component {
	c.Header("Content-Type", "text/html")
	c.Status(200)
	return pages.ImageForm(id, captchaBase64)
}

func (s *Site) FileForm(c *gin.Context, id, captchaBase64 string) templ.Component {
	c.Header("Content-Type", "text/html")
	c.Status(200)
	return pages.FileForm(id, captchaBase64)
}

func (s *Site) LinkResult(c *gin.Context, url string, qrBase64 string) templ.Component {
	c.Header("Content-Type", "text/html")
	c.Status(200)
	return pages.LinkResult(pages.LinkResultParams{Shortlink: url, QRBase64: qrBase64})
}

func (s *Site) ImageView(c *gin.Context, p pages.ImageViewParams) templ.Component {
	c.Header("Content-Type", "text/html")
	c.Status(200)
	return pages.ImageView(p)
}

func (s *Site) FileView(c *gin.Context, p pages.FileViewParams) templ.Component {
	c.Header("Content-Type", "text/html")
	c.Status(200)
	return pages.FileView(p)
}

func (s *Site) FileDownload(c *gin.Context, fileRawUrl string) templ.Component {
	c.Header("Content-Type", "text/html")
	c.Status(200)
	return pages.FileDownload(fileRawUrl)
}

func (s *Site) UserLogin(c *gin.Context) templ.Component {
	c.Header("Content-Type", "text/html")
	c.Status(200)
	return pages.LoginForm()
}

func (s *Site) UserRegister(c *gin.Context) templ.Component {
	c.Header("Content-Type", "text/html")
	c.Status(200)
	return pages.RegisterForm()
}

func (s *Site) AccountView(c *gin.Context, p pages.AccountViewParams) templ.Component {
	c.Header("Content-Type", "text/html")
	c.Status(200)
	return pages.AccountViewLinks(p)
}
