package server

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Base code got from gin source
// Needed to add ETag support
func StaticFS(router *gin.Engine, relativePath string, fs http.FileSystem) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}

	absolutePath := joinPaths(router.BasePath(), relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	etagMap := sync.Map{}

	handler := func(c *gin.Context) {
		file := c.Param("filepath")
		if etag, ok := etagMap.Load(file); ok {
			oldEtag := c.GetHeader("If-None-Match")
			if ok && oldEtag == etag.(string) {
				c.AbortWithStatus(http.StatusNotModified)
				return
			}
		}

		if _, noListing := fs.(*onlyFilesFS); noListing {
			c.Writer.WriteHeader(http.StatusNotFound)
		}

		f, err := fs.Open(file)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		fileBytes, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		hashBytes := md5.Sum(fileBytes)
		hash := hex.EncodeToString(hashBytes[:])
		etag := fmt.Sprintf("W/%s", hash)
		etagMap.Store(file, etag)

		c.Header("ETag", etag)
		c.Header("Cache-Control", "public, max-age=300")

		// TODO: no open to file twice, pass fileBytes here
		fileServer.ServeHTTP(c.Writer, c.Request)
	}

	urlPattern := path.Join(relativePath, "/*filepath")
	router.GET(urlPattern, handler)
	router.HEAD(urlPattern, handler)
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

type onlyFilesFS struct {
	fs http.FileSystem
}

func (fs onlyFilesFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(_ int) ([]os.FileInfo, error) {
	// this disables directory listing
	return nil, nil
}
