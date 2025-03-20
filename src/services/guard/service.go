package guard

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"shorty/src/common"
	"shorty/src/common/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrTemporaryBanned = errors.New("temporary banned")
	ErrTooManyRequests = errors.New("too many requests")
	ErrInternal        = errors.New("internal error")

	ErrWrongCaptcha  = errors.New("wrong captcha")
	ErrNoSuchCaptcha = errors.New("no such captcha")

	captchaSecret = common.NewShortId(10)
	tokenSecret   = common.NewShortId(10)
)

const (
	LimitWindow = 1 * time.Minute
	LimitAmount = 60

	BanWindow = 1 * time.Hour
	BanAmount = 2 * LimitAmount

	CaptchaTTL = 2 * time.Minute
)

func NewService(rdb *redis.Client, log logger.Logger, tracer trace.Tracer) *Service {
	return &Service{
		log:     log.WithService("guard"),
		tracer:  tracer,
		storage: newStorage(rdb, tracer),
	}
}

type Service struct {
	log     logger.Logger
	tracer  trace.Tracer
	storage *storage
}

func (s *Service) CheckIP(ctx context.Context, ip string) error {
	log := s.log.WithContext(ctx).WithService("guard")

	ctx, span := s.tracer.Start(ctx, "guard::CheckIP")
	defer span.End()

	banned, err := s.storage.IsBanned(ctx, ip)
	if err != nil {
		log.Error().Err(err).Msgf("checking banned with storage")
		return ErrInternal
	}
	if banned {
		log.Info().Msgf("rejected banned %s", ip)
		return ErrTemporaryBanned
	}

	rate, err := s.storage.IncRate(ctx, ip, LimitWindow)
	if err != nil {
		log.Error().Err(err).Msgf("inc requests rate with storage")
		return ErrInternal
	}
	if rate >= BanAmount {
		if err := s.storage.SetBanned(ctx, ip, BanWindow); err != nil {
			log.Error().Err(err).Msgf("set banned with storage")
			return ErrInternal
		}

		log.Info().Msgf("temporary banned %s", ip)
		return ErrTemporaryBanned
	}
	if rate >= LimitAmount {
		log.Info().Msgf("too many requests from %s", ip)
		return ErrTooManyRequests
	}

	return nil
}

func (s *Service) hashsum(value string) string {
	hashBytes := sha1.Sum([]byte(value + captchaSecret))
	return hex.EncodeToString(hashBytes[:])
}

func (s *Service) CreateCaptcha(ctx context.Context) (*CaptchaDTO, error) {
	log := s.log.WithContext(ctx).WithService("guard")

	ctx, span := s.tracer.Start(ctx, "guard::CreateCaptcha")
	defer span.End()

	id := common.NewShortId(16)
	value := common.NewDigitsString(4)

	hash := s.hashsum(value + captchaSecret)
	s.storage.SetCaptchaHash(ctx, id, hash, CaptchaTTL)

	imageBase64 := common.NewCaptchaImageBase64(id, value)

	log.Info().Msgf("created captcha, id=%s, value=%s", common.MaskSecret(id), common.MaskSecret(value))

	return &CaptchaDTO{id, imageBase64}, nil
}

func (s *Service) CheckCaptcha(ctx context.Context, id, value string) error {
	log := s.log.WithContext(ctx).WithService("guard")

	ctx, span := s.tracer.Start(ctx, "guard::CheckCaptcha")
	defer span.End()

	storedHash, err := s.storage.PopCaptchaHash(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting captcha hash id=%s", id)
		return ErrInternal
	}
	if storedHash == "" {
		log.Info().Msgf("no such captcha with id=%s", common.MaskSecret(id), common.MaskSecret(value))
		return ErrNoSuchCaptcha
	}

	if s.hashsum(value+captchaSecret) != storedHash {
		log.Info().Msgf("wrong captcha, id=%s, value=%s", common.MaskSecret(id), common.MaskSecret(value))
		return ErrWrongCaptcha
	}

	log.Info().Msgf("approved captcha, id=%s, value=%s", common.MaskSecret(id), common.MaskSecret(value))
	return nil
}

func (s *Service) CreateResourceToken(resource string, expiresAt int64) string {
	value := fmt.Sprintf("%s%d%s", resource, expiresAt, tokenSecret)
	return s.hashsum(value)
}

func (s *Service) CheckResourceToken(resource string, expiresAt int64, token string) bool {
	value := fmt.Sprintf("%s%d%s", resource, expiresAt, tokenSecret)
	return s.hashsum(value) == token
}
