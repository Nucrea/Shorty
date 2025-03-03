package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func NewService(baseUrl string, storage *Storage) *Service {
	return &Service{
		baseUrl: baseUrl,
		storage: storage,
	}
}

type Service struct {
	baseUrl string
	storage *Storage
}

func (s *Service) CreateLink(ctx context.Context, url string) (string, error) {
	linkId, err := s.storage.Create(ctx, url)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.baseUrl, linkId), nil
}

func (s *Service) CreateQR(ctx context.Context, url string) (string, error) {
	linkId, err := s.storage.Create(ctx, url)
	if err != nil {
		return "", err
	}
	link := fmt.Sprintf("%s/%s", s.baseUrl, linkId)

	qrc, err := qrcode.New(link)
	if err != nil {
		return "", err
	}

	buff := bytes.NewBuffer(nil)
	wrapper := BytesWrapper{buff}

	wr := standard.NewWithWriter(wrapper, standard.WithQRWidth(40))
	if err := qrc.Save(wr); err != nil {
		return "", err
	}

	result := base64.StdEncoding.EncodeToString(buff.Bytes())
	return result, nil
}

type BytesWrapper struct {
	io.Writer
}

func (BytesWrapper) Close() error {
	return nil
}
