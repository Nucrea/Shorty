package server

import (
	"shorty/src/common"

	"github.com/gin-gonic/gin"
)

type PostLinkInput struct {
	Url string `json:"url"`
}

type LinkOutput struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func (s *server) PostLink(ctx *gin.Context, input PostLinkInput) (*LinkOutput, error) {
	url, err := common.ValidateUrl(input.Url)
	if err != nil {
		return nil, err
	}

	link, err := s.LinksService.Save(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	return &LinkOutput{Id: link.Id, Url: url}, nil
}

func (s *server) GetLink(ctx *gin.Context, input struct{}) (*LinkOutput, error) {
	id, err := common.ValidateUrl(ctx.Param("id"))
	if err != nil {
		return nil, &ErrorBadRequest{err.Error()}
	}

	link, err := s.LinksService.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &LinkOutput{Id: link.Id, Url: link.Url}, nil
}
