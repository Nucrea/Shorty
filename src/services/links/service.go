package links

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
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

func NewService(pgConn *pgx.Conn, log *zerolog.Logger, appUrl string, tracer trace.Tracer) *Service {
	shortIdRegexp := regexp.MustCompile(`^\w{10}$`)
	newLog := log.With().Str("service", "links").Logger()
	return &Service{
		log:           &newLog,
		tracer:        tracer,
		shortIdRegexp: shortIdRegexp,
		appUrl:        appUrl,
		storage:       NewStorage(pgConn, tracer),
	}
}

type Service struct {
	log           *zerolog.Logger
	tracer        trace.Tracer
	shortIdRegexp *regexp.Regexp
	appUrl        string
	storage       *storage
}

func (s *Service) GetByShortId(ctx context.Context, linkId string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "links::GetByShortId")
	defer span.End()

	if !s.shortIdRegexp.MatchString(linkId) {
		return "", ErrBadShortId
	}

	link, err := s.storage.GetLink(ctx, linkId)
	if err != nil {
		s.log.Error().Err(err).Msgf("getting link with id=%s from storage", linkId)
		return "", ErrInternal
	}
	if link == "" {
		s.log.Info().Msgf("no such link with id=%s", linkId)
		return "", ErrNoSuchLink
	}
	return link, err
}

func (s *Service) CreateLink(ctx context.Context, url string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "links::CreateLink")
	defer span.End()

	if len(url) > 2000 {
		s.log.Info().Msgf("too long input url %s", url)
		return "", ErrBadUrl
	}

	url = strings.TrimSpace(url)
	if !govalidator.IsURL(url) {
		s.log.Info().Msgf("invalid input url %s", url)
		return "", ErrBadUrl
	}

	if !strings.HasSuffix(url, "http") || !strings.HasSuffix(url, "https") {
		url = fmt.Sprintf("https://%s", url)
	}

	shortId := NewShortId(10)
	_, err := s.storage.CreateLink(ctx, shortId, url)
	if err != nil {
		s.log.Error().Err(err).Msgf("creating link with storage")
		return "", ErrInternal
	}

	s.log.Info().Msgf("created link with id=%s", shortId)

	return fmt.Sprintf("%s/%s", s.appUrl, shortId), nil
}

func (s *Service) CreateQR(ctx context.Context, url string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "links::CreateQR")
	defer span.End()

	link, err := s.CreateLink(ctx, url)
	if err != nil {
		return "", err
	}

	qrc, err := qrcode.New(link)
	if err != nil {
		s.log.Error().Msgf("creating QR Code for link=%s", url)
		return "", ErrInternal
	}

	buf := bytes.NewBuffer(nil)
	wc := BytesWriterCloser{buf}

	wr := standard.NewWithWriter(wc, standard.WithQRWidth(40))
	if err := qrc.Save(wr); err != nil {
		s.log.Error().Msgf("saving QR Code for link=%s", url)
		return "", ErrInternal
	}

	s.log.Info().Msgf("created QR Code for link=%s", url)

	result := base64.StdEncoding.EncodeToString(buf.Bytes())
	return result, nil
}
