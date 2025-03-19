package guard

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"shorty/src/common"
	"shorty/src/common/cache"
	"shorty/src/common/logger"
	"time"

	"github.com/dchest/captcha"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrTemporaryBanned = errors.New("temporary banned")
	ErrTooManyRequests = errors.New("too many requests")
	ErrInternal        = errors.New("internal error")
	ErrWrongCaptcha    = errors.New("wrong captcha")

	CaptchaSecret = common.NewShortId(10)
)

const (
	LimitWindow = 1 * time.Minute
	LimitAmount = 60

	BanWindow = 1 * time.Hour
	BanAmount = 2 * LimitAmount
)

func NewService(rdb *redis.Client, log logger.Logger, tracer trace.Tracer) *Service {
	return &Service{
		log:          log.WithService("ratelimit"),
		tracer:       tracer,
		storage:      newStorage(rdb, tracer),
		captchaCache: cache.NewInmem[string](),
	}
}

type Service struct {
	log          logger.Logger
	tracer       trace.Tracer
	storage      *storage
	captchaCache *cache.Inmem[string]
}

func (s *Service) CheckIP(ctx context.Context, ip string) error {
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "ratelimit::Check")
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

func (s *Service) CreateCaptcha() (*CaptchaDTO, error) {
	id := common.NewShortId(16)
	value := common.NewDigitsString(4)

	digits := make([]byte, 4)
	for i, c := range value {
		digits[i] = byte(c - '0')
	}

	buf := bytes.NewBuffer(nil)
	image := captcha.NewImage(id, digits, 200, 80)
	image.WriteTo(buf)

	hashBytes := sha1.Sum([]byte(value + CaptchaSecret))
	hash := hex.EncodeToString(hashBytes[:])
	s.captchaCache.SetEx(id, hash, time.Minute)

	return &CaptchaDTO{
		Id:          id,
		ImageBase64: base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}

func (s *Service) CheckCaptcha(id, captcha string) error {
	storedHash, ok := s.captchaCache.Get(id)
	if !ok {
		return ErrWrongCaptcha
	}

	hashBytes := sha1.Sum([]byte(captcha + CaptchaSecret))
	hash := hex.EncodeToString(hashBytes[:])

	if hash != storedHash {
		return ErrWrongCaptcha
	}

	s.captchaCache.Del(id)
	return nil
}
