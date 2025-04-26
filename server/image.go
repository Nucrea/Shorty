package server

import (
	"fmt"
	"io"
	"net/http"
	"shorty/src/common"
	"shorty/src/services/image"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ImageInfoOutput struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	OriginalUrl  string `json:"originalUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

func (s *server) UploadImage(c *gin.Context, input DummyInput) (*ImageInfoOutput, error) {
	header, err := c.FormFile("file")
	if err != nil {
		s.Logger.Error().Err(err).Msg("aaaaaaa")
		return nil, &ErrorBadRequest{err.Error()}
	}

	file, err := header.Open()
	if err != nil {
		log.Error().Err(err).Msg("error opening tmp file")
		return nil, &ErrorInternal{err.Error()}
	}
	defer file.Close()

	fileBytes, _ := io.ReadAll(file)
	image, err := s.ImageService.UploadImage(c, header.Filename, fileBytes)
	if err != nil {
		return nil, err
	}

	originalUrl := fmt.Sprintf("%s/i/o/%s", s.Url, image.Id)
	thumbnailUrl := fmt.Sprintf("%s/i/o/%s", s.Url, image.Id)

	return &ImageInfoOutput{
		Id:           image.Id,
		Name:         image.Name,
		Size:         10,
		OriginalUrl:  originalUrl,
		ThumbnailUrl: thumbnailUrl,
	}, nil
}

func (s *server) GetImageInfo(ctx *gin.Context, input DummyInput) (*ImageInfoOutput, error) {
	id := ctx.Param("id")
	if !common.ValidateShortId(id) {
		return nil, &ErrorBadRequest{"bad id"}
	}

	image, err := s.ImageService.GetImageMetadata(ctx, id)
	if err != nil {
		return nil, err
	}

	originalUrl := fmt.Sprintf("%s/i/o/%s", s.Url, image.Id)
	thumbnailUrl := fmt.Sprintf("%s/i/o/%s", s.Url, image.Id)

	return &ImageInfoOutput{
		Id:           image.Id,
		Name:         image.Name,
		Size:         10,
		OriginalUrl:  originalUrl,
		ThumbnailUrl: thumbnailUrl,
	}, nil
}

func (s *server) ResolveImage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(400)
		return
	}

	isThumbnail := false
	switch c.Param("type") {
	case "o":
		isThumbnail = false
	case "t":
		isThumbnail = true
	default:
		c.AbortWithStatus(404)
		return
	}

	meta, err := s.ImageService.GetImageMetadata(c, id)
	if err == image.ErrImageNotFound {
		c.AbortWithStatus(404)
		return
	}
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	if oldEtag := c.GetHeader("If-None-Match"); oldEtag == meta.Hash {
		s.Logger.WithContext(c).Info().Msgf("image (id=%s) hash does not changed", id)
		c.AbortWithStatus(http.StatusNotModified)
		return
	}

	imageBytes, err := s.ImageService.GetImageBytes(c, id, isThumbnail)
	if err == image.ErrImageNotFound {
		c.AbortWithStatus(404)
		return
	}
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.Header("ETag", meta.Hash)
	c.Header("Cache-Control", "public, max-age=300")
	c.Data(200, "image/jpeg", imageBytes)
}
