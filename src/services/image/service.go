package image

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"shorty/src/common/broker"
	"shorty/src/common/logger"

	"github.com/anthonynsimon/bild/transform"
	"github.com/jackc/pgx/v5/pgxpool"
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

func NewService(pg *pgxpool.Pool, s3 *minio.Client, log logger.Logger, tracer trace.Tracer) *Service {
	return &Service{
		log:         log.WithService("image"),
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

func (s *Service) createThumbnail(ctx context.Context, imgBytes []byte) ([]byte, error) {
	log := s.log.WithContext(ctx)

	_, span := s.tracer.Start(ctx, "image::createThumbnail")
	defer span.End()

	imgReader := bytes.NewReader(imgBytes)
	imgInfo, format, err := image.DecodeConfig(imgReader)
	if err != nil {
		log.Error().Err(err).Msg("failed decoding image")
		return nil, ErrInvalidFormat
	}
	if format != "jpeg" { //&& format != "png" {
		log.Info().Msgf("rejected image with unsupported format %s", format)
		return nil, ErrUnsupportedFormat
	}

	imgReader.Seek(0, io.SeekStart)
	img, _, err := image.Decode(imgReader)
	if err != nil {
		log.Error().Err(err).Msg("failed decoding image")
		return nil, ErrInvalidFormat
	}

	width := 200
	height := width * imgInfo.Height / imgInfo.Width

	resized := transform.Resize(img, width, height, transform.Linear)
	buff := bytes.NewBuffer(nil)
	if err := jpeg.Encode(buff, resized, nil); err != nil {
		log.Error().Err(err).Msg("failed encoding thumbnail")
		return nil, ErrInternal
	}

	return buff.Bytes(), nil
}

func (s *Service) getHash(ctx context.Context, imgBytes []byte) string {
	_, span := s.tracer.Start(ctx, "image::getHash")
	defer span.End()

	hash := sha512.Sum512(imgBytes)
	return hex.EncodeToString(hash[:])
}

func (s *Service) UploadImage(ctx context.Context, name string, imgBytes []byte) (*ImageInfoDTO, error) {
	log := s.log.WithContext(ctx)

	_, span := s.tracer.Start(ctx, "image::UploadImage")
	defer span.End()

	if size := len(imgBytes); size > 5*1024*1024 { //temporary 15MB max
		log.Info().Msgf("rejected too heavy image with size %d", size)
		return nil, ErrImageTooLarge
	}

	hash := s.getHash(ctx, imgBytes)
	info, err := s.infoStorage.GetImageInfoByHash(ctx, hash)
	if err != nil {
		log.Error().Err(err).Msg("failed getting img info by hash")
		return nil, ErrInternal
	}

	var (
		foundExisting = info != nil
		imgSize       = len(imgBytes)
		imageId       string
		thumbId       string
	)
	if foundExisting {
		log.Info().Msg("found existing files with same hash, add reference to them")
		imageId = info.ImageId
		thumbId = info.ThumbnailId
		imgSize = info.Size
	} else {
		log.Info().Msg("not found existing files with same hash, saving img and thumb to storage...")

		thumbBytes, err := s.createThumbnail(ctx, imgBytes)
		if err != nil {
			return nil, err
		}

		imageId = NewShortId(32)
		if err := s.fileStorage.SaveFile(ctx, imageId, imgBytes); err != nil {
			log.Error().Err(err).Msg("failed saving main file")
			return nil, ErrInternal
		}

		thumbId = NewShortId(32)
		if err := s.fileStorage.SaveFile(ctx, thumbId, thumbBytes); err != nil {
			// s.broker.PutFilesToDelete(ctx, imageId)
			log.Error().Err(err).Msg("failed saving thumb file")
			return nil, ErrInternal
		}
	}

	shortId := NewShortId(32)
	result, err := s.infoStorage.SaveImageInfo(ctx,
		ImageInfoDTO{
			ShortId:     shortId,
			Size:        imgSize,
			Hash:        hash,
			Name:        name,
			ImageId:     imageId,
			ThumbnailId: thumbId,
		},
	)
	if err != nil {
		// s.broker.PutFilesToDelete(ctx, imageId, thumbId)
		log.Error().Err(err).Msg("failed writing info")
		return nil, ErrInternal
	}

	log.Info().Msgf("created image with id=%s", result.ShortId)

	return result, nil
}

func (s *Service) GetImageInfo(ctx context.Context, shortId string) (*ImageInfoDTO, error) {
	log := s.log.WithContext(ctx)

	_, span := s.tracer.Start(ctx, "image::GetImageInfo")
	defer span.End()

	imageInfo, err := s.infoStorage.GetImageInfoByShortId(ctx, shortId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting image (id=%s) info from storage", shortId)
		return nil, ErrInternal
	}
	if imageInfo == nil {
		log.Info().Msgf("not found image with id=%s", shortId)
		return nil, ErrImageNotFound
	}

	return imageInfo, nil
}

func (s *Service) GetFile(ctx context.Context, shortId string) ([]byte, error) {
	log := s.log.WithContext(ctx)

	_, span := s.tracer.Start(ctx, "image::GetFile")
	defer span.End()

	// imageInfo, err := s.infoStorage.GetImageInfo(ctx, shortId)
	// if err != nil {
	// 	log.Error().Err(err).Msgf("failed getting image (id=%s) info from storage", shortId)
	// 	return nil, ErrInternal
	// }
	// if imageInfo == nil {
	// 	log.Info().Msgf("not found image with id=%s", shortId)
	// 	return nil, ErrImageNotFound
	// }

	imgBytes, err := s.fileStorage.GetFile(ctx, shortId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting image (id=%s) from storage", shortId)
		return nil, ErrInternal
	}

	return imgBytes, nil
}

func (s *Service) GetImage(ctx context.Context, shortId string) (*ImageDTO, error) {
	log := s.log.WithContext(ctx)

	_, span := s.tracer.Start(ctx, "image::GetImage")
	defer span.End()

	imageInfo, err := s.infoStorage.GetImageInfoByShortId(ctx, shortId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting image (id=%s) info from storage", shortId)
		return nil, ErrInternal
	}
	if imageInfo == nil {
		log.Info().Msgf("not found image with id=%s", shortId)
		return nil, ErrImageNotFound
	}

	imageBytes, err := s.fileStorage.GetFile(ctx, imageInfo.ImageId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting image (id=%s) from storage", shortId)
		return nil, ErrInternal
	}

	return &ImageDTO{
		Id:    imageInfo.ImageId,
		Size:  imageInfo.Size,
		Name:  imageInfo.Name,
		Bytes: imageBytes,
	}, nil
}
