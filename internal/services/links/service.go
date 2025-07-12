package links

import (
	"context"
	"errors"
	"shorty/internal/common"
	"shorty/internal/common/logging"
	"shorty/internal/common/metrics"

	"go.opentelemetry.io/otel/trace"
)

var (
	ErrBadShortId = errors.New("invalid short id")
	ErrBadUrl     = errors.New("invalid url")
	ErrNoSuchLink = errors.New("no such link")
	ErrInternal   = errors.New("internal error")
)

func NewService(storage Storage, logger logging.Logger, tracer trace.Tracer, meter metrics.Meter) *Service {
	return &Service{
		logger:          logger.WithService("links"),
		tracer:          tracer,
		storage:         storage,
		createdCounter:  meter.NewCounter("links_created", "Created links counter"),
		resolvedCounter: meter.NewCounter("links_resolved", "Resolved links counter"),
	}
}

type Service struct {
	logger  logging.Logger
	tracer  trace.Tracer
	storage Storage

	createdCounter  metrics.Counter
	resolvedCounter metrics.Counter
}

func (s *Service) GetByShortId(ctx context.Context, linkId string) (string, error) {
	log := s.logger.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::GetByShortId")
	defer span.End()

	if !common.ValidateShortId(linkId) {
		return "", ErrBadShortId
	}

	link, err := s.storage.GetShortlink(ctx, linkId)
	if err != nil {
		log.Error().Err(err).Msgf("getting link with id=%s from storage", linkId)
		return "", ErrInternal
	}
	if link == "" {
		log.Info().Msgf("no such link with id=%s", linkId)
		return "", ErrNoSuchLink
	}

	log.Info().Msgf("got link with id=%s from storage", linkId)
	s.resolvedCounter.Inc()

	return link, nil
}

func (s *Service) Create(ctx context.Context, url string) (string, error) {
	log := s.logger.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::CreateShortlink")
	defer span.End()

	url = common.ValidateUrl(url)
	if url == "" {
		log.Info().Msgf("invalid input url %s", url)
		return "", ErrBadUrl
	}

	id := common.NewShortId(10)
	if err := s.storage.SaveShortlink(ctx, id, url); err != nil {
		log.Error().Err(err).Msgf("creating qr and link with storage")
		return "", ErrInternal
	}

	log.Info().Msgf("created shortlink with id=%s", id)
	s.createdCounter.Inc()

	return id, nil
}
