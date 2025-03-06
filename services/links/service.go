package links

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgx/v5"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

var (
	ErrBadShortId = fmt.Errorf("invalid short id")
	ErrBadUrl     = fmt.Errorf("invalid url")
	ErrNoSuchLink = fmt.Errorf("no such link")
)

func NewService(pgConn *pgx.Conn, baseUrl string) *Service {
	shortIdRegexp := regexp.MustCompile(`^\w{10}$`)
	return &Service{
		shortIdRegexp: shortIdRegexp,
		baseUrl:       baseUrl,
		storage:       NewStorage(pgConn),
	}
}

type Service struct {
	shortIdRegexp *regexp.Regexp
	baseUrl       string
	storage       *Storage
}

func (s *Service) GetByShortId(ctx context.Context, linkId string) (string, error) {
	if !s.shortIdRegexp.MatchString(linkId) {
		return "", ErrBadShortId
	}

	link, err := s.storage.Get(ctx, linkId)
	if err != nil {
		return "", err
	}
	if link == "" {
		return "", ErrNoSuchLink
	}
	return link, err
}

func (s *Service) CreateLink(ctx context.Context, url string) (string, error) {
	if len(url) > 2000 {
		return "", ErrBadUrl
	}

	url = strings.TrimSpace(url)
	if !govalidator.IsURL(url) {
		return "", ErrBadUrl
	}

	if !strings.HasSuffix(url, "http") || !strings.HasSuffix(url, "https") {
		url = fmt.Sprintf("https://%s", url)
	}

	shortId := NewShortId(10)
	_, err := s.storage.Create(ctx, shortId, url)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", s.baseUrl, shortId), nil
}

func (s *Service) CreateQR(ctx context.Context, url string) (string, error) {
	link, err := s.CreateLink(ctx, url)
	if err != nil {
		return "", err
	}

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
