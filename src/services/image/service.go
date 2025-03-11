package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"shorty/src/common/broker"
	"shorty/src/common/logger"

	"github.com/anthonynsimon/bild/transform"
	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrUnsupportedFormat = fmt.Errorf("unsupported format")
	ErrImageNotFound     = fmt.Errorf("image not found")
	ErrImageTooLarge     = fmt.Errorf("image too large")
	ErrInvalidFormat     = fmt.Errorf("invalid format")
	ErrInternal          = fmt.Errorf("internal error")
)

func NewService(pg *pgx.Conn, s3 *minio.Client, log logger.Logger, tracer trace.Tracer) *Service {
	return &Service{
		log:         log,
		tracer:      tracer,
		fileStorage: newFileStorage(s3, tracer),
		infoStorage: newInfoStorage(pg, tracer),
	}
}

type Service struct {
	log         logger.Logger
	tracer      trace.Tracer
	broker      broker.Broker
	fileStorage *fileStorage
	infoStorage *infoStorage
}

func (s *Service) CreateImage(ctx context.Context, name string, imgBytes []byte) (*ImageInfoDTO, error) {
	log := s.log.WithContext(ctx)

	_, span := s.tracer.Start(ctx, "image::CreateImage")
	defer span.End()

	if size := len(imgBytes); size > 15*1024*1024 { //temporary 15MB max
		log.Info().Msgf("rejected too heavy image with size %d", size)
		return nil, ErrImageTooLarge
	}

	img, format, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		log.Error().Err(err).Msg("failed decoding image")
		return nil, ErrInvalidFormat
	}
	if format != "jpeg" { //&& format != "png" {
		log.Info().Msgf("rejected image with unsupported format %s", format)
		return nil, ErrUnsupportedFormat
	}

	resized := transform.Resize(img, 600, 400, transform.Linear)
	buff := bytes.NewBuffer(nil)
	if err := jpeg.Encode(buff, resized, nil); err != nil {
		log.Error().Err(err).Msg("failed encoding thumbnail")
		return nil, ErrInternal
	}
	thumbBytes := buff.Bytes()

	imageId := NewShortId(32)
	if err := s.fileStorage.SaveFile(ctx, imageId, imgBytes); err != nil {
		log.Error().Err(err).Msg("failed saving main file")
		return nil, ErrInternal
	}

	thumbId := NewShortId(32)
	if err := s.fileStorage.SaveFile(ctx, thumbId, thumbBytes); err != nil {
		s.broker.PutFilesToDelete(ctx, imageId)
		log.Error().Err(err).Msg("failed saving thumb file")
		return nil, ErrInternal
	}

	shortId := NewShortId(32)
	result, err := s.infoStorage.SaveImageInfo(ctx,
		ImageInfoDTO{
			ShortId:     shortId,
			Size:        len(imgBytes),
			Name:        name,
			ImageId:     imageId,
			ThumbnailId: thumbId,
		},
	)
	if err != nil {
		s.broker.PutFilesToDelete(ctx, imageId, thumbId)
		log.Error().Err(err).Msg("failed writing info")
		return nil, ErrInternal
	}

	log.Info().Msgf("created image with id=%s", result.ShortId)

	return result, nil
}

func (s *Service) GetThumbnail(ctx context.Context, shortId string) ([]byte, error) {
	log := s.log.WithContext(ctx)

	_, span := s.tracer.Start(ctx, "image::GetThumbnail")
	defer span.End()

	imageInfo, err := s.infoStorage.GetImageInfo(ctx, shortId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting image (id=%s) info from storage", shortId)
		return nil, ErrInternal
	}

	thumbBytes, err := s.fileStorage.GetFile(ctx, imageInfo.ThumbnailId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting thumbnail (id=%s) from storage", shortId)
		return nil, ErrInternal
	}

	return thumbBytes, nil
}

func (s *Service) GetImage(ctx context.Context, shortId string) ([]byte, error) {
	log := s.log.WithContext(ctx)

	_, span := s.tracer.Start(ctx, "image::GetImage")
	defer span.End()

	imageInfo, err := s.infoStorage.GetImageInfo(ctx, shortId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting image (id=%s) info from storage", shortId)
		return nil, ErrInternal
	}

	imageBytes, err := s.fileStorage.GetFile(ctx, imageInfo.ImageId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting image (id=%s) from storage", shortId)
		return nil, ErrInternal
	}

	return imageBytes, nil
}
