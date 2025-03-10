package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
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
	fileStorage *fileStorage
	infoStorage *infoStorage
}

func (s *Service) CreateImage(ctx context.Context, name string, imgBytes []byte) (*ImageInfoDTO, error) {
	if len(imgBytes) > 15*1024*1024 { //temporary 15MB max
		return nil, ErrImageTooLarge
	}

	img, format, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, ErrInvalidFormat
	}
	if format != "jpeg" { //&& format != "png" {
		return nil, ErrUnsupportedFormat
	}

	resized := transform.Resize(img, 800, 600, transform.Linear)
	buff := bytes.NewBuffer(nil)
	if err := jpeg.Encode(buff, resized, nil); err != nil {
		return nil, ErrInternal
	}
	thumbBytes := buff.Bytes()

	imageId := NewShortId(32)
	if err := s.fileStorage.SaveFile(ctx, imageId, imgBytes); err != nil {
		s.log.Error().Err(err).Msg("failed saving main file")
		return nil, ErrInternal
	}

	thumbId := NewShortId(32)
	if err := s.fileStorage.SaveFile(ctx, thumbId, thumbBytes); err != nil {
		s.log.Error().Err(err).Msg("failed saving thumb file")
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
		s.log.Error().Err(err).Msg("failed writing info")
		return nil, ErrInternal
	}

	return result, nil
}

func (s *Service) GetThumbnail(ctx context.Context, shortId string) ([]byte, error) {
	imageInfo, err := s.infoStorage.GetImageInfo(ctx, shortId)
	if err != nil {
		return nil, ErrInternal
	}

	thumbBytes, err := s.fileStorage.GetFile(ctx, imageInfo.ThumbnailId)
	if err != nil {
		return nil, ErrInternal
	}

	return thumbBytes, nil
}

func (s *Service) GetImage(ctx context.Context, shortId string) ([]byte, error) {
	imageInfo, err := s.infoStorage.GetImageInfo(ctx, shortId)
	if err != nil {
		return nil, ErrInternal
	}

	imageBytes, err := s.fileStorage.GetFile(ctx, imageInfo.ImageId)
	if err != nil {
		return nil, ErrInternal
	}

	return imageBytes, nil
}
