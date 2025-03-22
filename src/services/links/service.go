package links

import (
	"context"
	"errors"
	"shorty/src/common"
	"shorty/src/common/logging"
	"shorty/src/common/metrics"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrBadShortId = errors.New("invalid short id")
	ErrBadUrl     = errors.New("invalid url")
	ErrNoSuchLink = errors.New("no such link")
	ErrInternal   = errors.New("internal error")
)

func NewService(pgConn *pgxpool.Pool, log logging.Logger, appUrl string, tracer trace.Tracer, meter metrics.Meter) *Service {
	return &Service{
		log:     log.WithService("links"),
		tracer:  tracer,
		storage: NewStorage(pgConn, tracer, meter),
	}
}

type Service struct {
	log     logging.Logger
	tracer  trace.Tracer
	storage *storage
}

func (s *Service) GetByShortId(ctx context.Context, linkId string) (string, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::GetByShortId")
	defer span.End()

	if !common.ValidateShortId(linkId) {
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

func (s *Service) Create(ctx context.Context, url string) (string, error) {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::CreateShortlink")
	defer span.End()

	url = common.ValidateUrl(url)
	if url == "" {
		log.Info().Msgf("invalid input url %s", url)
		return "", ErrBadUrl
	}

	id := common.NewShortId(10)
	if err := s.storage.CreateLink(ctx, id, url); err != nil {
		log.Error().Err(err).Msgf("creating qr and link with storage")
		return "", ErrInternal
	}

	log.Info().Msgf("created shortlink with id=%s", id)

	return id, nil
}
