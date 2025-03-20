package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"shorty/src/common"
	"shorty/src/common/assets"
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

const (
	BucketName   = "images"
	MaxImageSize = 5 * 1024 * 1024
)

func NewService(pg *pgxpool.Pool, s3 *minio.Client, log logger.Logger, tracer trace.Tracer) *Service {
	return &Service{
		log:          log.WithService("image"),
		tracer:       tracer,
		assetStorage: assets.NewStorage(pg, s3, tracer, log),
		metaRepo:     newMetadataRepo(pg, tracer),
	}
}

type Service struct {
	log          logger.Logger
	tracer       trace.Tracer
	broker       broker.Broker
	assetStorage *assets.Storage
	metaRepo     *metadataRepo
}

func (s *Service) createThumbnail(ctx context.Context, imgBytes []byte) ([]byte, error) {
	log := s.log.WithContext(ctx).WithService("image")

	ctx, span := s.tracer.Start(ctx, "image::createThumbnail")
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

func (s *Service) UploadImage(ctx context.Context, name string, imageBytes []byte) (*ImageMetadataDTO, error) {
	log := s.log.WithContext(ctx).WithService("image")

	ctx, span := s.tracer.Start(ctx, "image::UploadImage")
	defer span.End()

	imageSize := len(imageBytes)
	if imageSize > MaxImageSize { //temporary 15MB max
		log.Info().Msgf("rejected too heavy image with size %d", imageSize)
		return nil, ErrImageTooLarge
	}

	imageHash := common.NewAssetHash(imageBytes)
	info, err := s.metaRepo.GetImageMetadataDuplicate(ctx, imageSize, imageHash)
	if err != nil {
		log.Error().Err(err).Msg("failed getting img info by hash")
		return nil, ErrInternal
	}

	metadata := ImageMetadataDTO{
		Id:   common.NewShortId(32),
		Name: name,
	}

	if info != nil {
		log.Info().Msg("found existing files with same hash, add reference to them")
		metadata.OriginalId = info.OriginalId
		metadata.ThumbnailId = info.ThumbnailId
	} else {
		log.Info().Msg("not found existing files with same hash, saving img and thumb to storage...")

		thumbBytes, err := s.createThumbnail(ctx, imageBytes)
		if err != nil {
			return nil, err
		}

		assets, err := s.assetStorage.SaveAssets(ctx, BucketName, imageBytes, thumbBytes)
		if err != nil {
			log.Error().Err(err).Msg("failed saving assets")
			return nil, ErrInternal
		}

		metadata.OriginalId = assets[0].Id
		metadata.ThumbnailId = assets[1].Id
	}

	err = s.metaRepo.SaveImageMetadata(ctx, metadata)
	if err != nil {
		log.Error().Err(err).Msg("failed saving image metadata")
		return nil, ErrInternal
	}

	log.Info().Msgf("created image with id=%s", metadata.Id)

	return &metadata, nil
}

func (s *Service) GetImageMetadata(ctx context.Context, id string) (*ImageMetadataExDTO, error) {
	log := s.log.WithContext(ctx).WithService("image")

	ctx, span := s.tracer.Start(ctx, "image::GetImageMetadata")
	defer span.End()

	meta, err := s.metaRepo.GetImageMetadataById(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting image (id=%s) info from storage", id)
		return nil, ErrInternal
	}
	if meta == nil {
		log.Info().Msgf("not found image with id=%s", id)
		return nil, ErrImageNotFound
	}

	log.Info().Msgf("read image metadata (id=%s)", id)

	return meta, nil
}

func (s *Service) GetImageBytes(ctx context.Context, id string, thumbnail bool) ([]byte, error) {
	log := s.log.WithContext(ctx).WithService("image")

	ctx, span := s.tracer.Start(ctx, "image::GetImageBytes")
	defer span.End()

	meta, err := s.GetImageMetadata(ctx, id)
	if err != nil {
		return nil, err
	}

	assetId := meta.OriginalId
	if thumbnail {
		assetId = meta.ThumbnailId
	}

	assetBytes, err := s.assetStorage.GetAssetBytes(ctx, BucketName, assetId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting image asset bytes from storage (id=%s, assetId=%s, thumbnail=%t)", id, assetId, thumbnail)
		return nil, ErrInternal
	}

	log.Info().Msgf("read image asset (id=%s, assetId=%s, thumbnail=%t)", id, assetId, thumbnail)

	return assetBytes, nil
}
