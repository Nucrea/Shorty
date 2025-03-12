package links

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"shorty/src/common/logger"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgx/v5"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrBadShortId = errors.New("invalid short id")
	ErrBadUrl     = errors.New("invalid url")
	ErrNoSuchLink = errors.New("no such link")
	ErrInternal   = errors.New("internal error")
)

func NewService(pgConn *pgx.Conn, log logger.Logger, appUrl string, tracer trace.Tracer) *Service {
	shortIdRegexp := regexp.MustCompile(`^\w{10}$`)
	return &Service{
		log:           log.WithService("links"),
		tracer:        tracer,
		shortIdRegexp: shortIdRegexp,
		appUrl:        appUrl,
		storage:       NewStorage(pgConn, tracer),
	}
}

type Service struct {
	log           logger.Logger
	tracer        trace.Tracer
	shortIdRegexp *regexp.Regexp
	appUrl        string
	storage       *storage
}

func (s *Service) GetByShortId(ctx context.Context, linkId string) (string, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::GetByShortId")
	defer span.End()

	if !s.shortIdRegexp.MatchString(linkId) {
		return "", ErrBadShortId
	}

	link, err := s.storage.GetLink(ctx, linkId)
	if err != nil {
		log.Error().Err(err).Msgf("getting link with id=%s from storage", linkId)
		return "", ErrInternal
	}
	if link == "" {
		log.Info().Msgf("no such link with id=%s", linkId)
		return "", ErrNoSuchLink
	}

	log.Info().Msgf("got link with id=%s from storage", linkId)

	return link, nil
}

func (s *Service) parseUrl(url string) (string, error) {
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

	return url, nil
}

func (s *Service) newLink(shortId string) string {
	return fmt.Sprintf("%s/s/%s", s.appUrl, shortId)
}

func (s *Service) CreateLink(ctx context.Context, url string) (string, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::CreateLink")
	defer span.End()

	url, err := s.parseUrl(url)
	if err != nil {
		log.Info().Msgf("invalid input url %s", url)
		return "", ErrBadUrl
	}

	shortId := NewShortId(10)
	if err := s.storage.CreateLink(ctx, shortId, url); err != nil {
		log.Error().Err(err).Msgf("creating link with storage")
		return "", ErrInternal
	}

	log.Info().Msgf("created link with id=%s", shortId)

	return s.newLink(shortId), nil
}

func (s *Service) CreateQR(ctx context.Context, url string) (string, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::CreateQR")
	defer span.End()

	url, err := s.parseUrl(url)
	if err != nil {
		log.Info().Msgf("invalid input url %s", url)
		return "", ErrBadUrl
	}

	shortId := NewShortId(10)
	link := s.newLink(shortId)

	qrc, err := qrcode.New(link)
	if err != nil {
		log.Error().Msgf("creating QR Code for link=%s", url)
		return "", ErrInternal
	}

	buf := bytes.NewBuffer(nil)
	wc := BytesWriterCloser{buf}

	wr := standard.NewWithWriter(wc, standard.WithQRWidth(40))
	if err := qrc.Save(wr); err != nil {
		log.Error().Msgf("saving QR Code for link=%s", url)
		return "", ErrInternal
	}

	if err := s.storage.CreateLink(ctx, shortId, url); err != nil {
		log.Error().Err(err).Msgf("creating qr and link with storage")
		return "", ErrInternal
	}

	log.Info().Msgf("created QR Code for link=%s", url)

	result := base64.StdEncoding.EncodeToString(buf.Bytes())
	return result, nil
}
