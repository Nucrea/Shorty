package links

import (
	"context"
	"errors"
	"shorty/src/common"
	"shorty/src/common/logging"
	"shorty/src/common/metrics"
	"strings"

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

func (s *Service) Save(ctx context.Context, url string, userId *string) (*LinkDTO, error) {
	log := s.logger.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::CreateShortlink")
	defer span.End()

	url = common.ValidateUrl(url)
	if url == "" {
		log.Info().Msgf("invalid input url %s", url)
		return nil, ErrBadUrl
	}
	id := common.NewShortId(10)

	var err error
	if userId != nil {
		err = s.storage.SaveLinkForUser(ctx, id, *userId, url)
	} else {
		err = s.storage.SaveLink(ctx, id, url)
	}
	if err != nil {
		log.Error().Err(err).Msgf("failed saving link to storage")
		return nil, ErrInternal
	}

	log.Info().Msgf("created shortlink with id=%s", id)
	s.createdCounter.Inc()

	return &LinkDTO{
		Id:  id,
		Url: url,
	}, nil
}

func (s *Service) Update(ctx context.Context, id, userId, url string) error {
	log := s.logger.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::CreateShortlink")
	defer span.End()

	url = common.ValidateUrl(url)
	if url == "" {
		log.Info().Msgf("invalid input url %s", url)
		return ErrBadUrl
	}
	if !common.ValidateShortId(id) {
		return ErrBadShortId
	}

	err := s.storage.SaveLinkForUser(ctx, id, userId, url)
	if err != nil {
		log.Error().Err(err).Msgf("failed saving link to storage")
		return ErrInternal
	}

	log.Info().Msgf("updated shortlink with id=%s", id)
	return nil
}

func (s *Service) GetById(ctx context.Context, linkId string) (*LinkDTO, error) {
	log := s.logger.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::GetByShortId")
	defer span.End()

	if !common.ValidateShortId(linkId) {
		return nil, ErrBadShortId
	}

	link, err := s.storage.GetLinkById(ctx, linkId)
	if err != nil {
		log.Error().Err(err).Msgf("getting link with id=%s from storage", linkId)
		return nil, ErrInternal
	}
	if link == nil {
		log.Info().Msgf("no such link with id=%s", linkId)
		return nil, ErrNoSuchLink
	}

	log.Info().Msgf("got link with id=%s from storage", linkId)
	s.resolvedCounter.Inc()

	return link, nil
}

func (s *Service) GetByUserId(ctx context.Context, userId string) ([]*LinkDTO, error) {
	log := s.logger.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::GetByUserId")
	defer span.End()

	links, err := s.storage.GetLinksByUserId(ctx, userId)
	if err != nil {
		log.Error().Err(err).Msgf("getting link for userId=%s from storage", userId)
		return nil, ErrInternal
	}
	if links == nil {
		log.Info().Msgf("no such links for userId=%s", userId)
		return nil, ErrNoSuchLink
	}

	log.Info().Msgf("got links with userId=%s from storage", userId)
	return links, nil
}

func (s *Service) Delete(ctx context.Context, ids ...string) error {
	log := s.logger.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "links::GetByUserId")
	defer span.End()

	idsStr := strings.Join(ids, ", ")

	err := s.storage.DeleteLinks(ctx, ids...)
	if err != nil {
		log.Error().Err(err).Msgf("failed deleting link for with id=[ %s ] from storage", idsStr)
		return ErrInternal
	}

	log.Info().Msgf("deleted link with ids=[ %s ] from storage", idsStr)
	return nil
}
