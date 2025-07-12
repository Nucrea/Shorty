package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"shorty/internal/common"
	"shorty/internal/common/broker"
	"shorty/internal/common/logging"
	"shorty/internal/common/metrics"
	"shorty/internal/services/assets"

	"github.com/anthonynsimon/bild/transform"
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

func NewService(metaRepo MetadataRepo, assetsStorage *assets.Storage, log logging.Logger, tracer trace.Tracer, meter metrics.Meter) *Service {
	return &Service{
		log:                   log.WithService("images"),
		tracer:                tracer,
		assetStorage:          assetsStorage,
		metaRepo:              metaRepo,
		uploadsCounter:        meter.NewCounter("images_uploads", "Count of uploaded images"),
		dulicatesCounter:      meter.NewCounter("images_duplicates", "Count of uploaded duplicates"),
		origDownloadsCounter:  meter.NewCounter("images_orig_downloads", "How many times original image was downloaded"),
		thumbDownloadsCounter: meter.NewCounter("images_thumb_downloads", "How many times thumbnail of image was downloaded"),
	}
}

type Service struct {
	log          logging.Logger
	tracer       trace.Tracer
	broker       broker.Broker
	assetStorage *assets.Storage
	metaRepo     MetadataRepo

	uploadsCounter        metrics.Counter
	dulicatesCounter      metrics.Counter
	origDownloadsCounter  metrics.Counter
	thumbDownloadsCounter metrics.Counter
}

func (s *Service) createThumbnail(ctx context.Context, imgBytes []byte) ([]byte, error) {
	log := s.log.WithContext(ctx)

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
	log := s.log.WithContext(ctx)

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
		s.dulicatesCounter.Inc()
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
	s.uploadsCounter.Inc()

	return &metadata, nil
}

func (s *Service) GetImageMetadata(ctx context.Context, id string) (*ImageMetadataExDTO, error) {
	log := s.log.WithContext(ctx)

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
	log := s.log.WithContext(ctx)

	ctx, span := s.tracer.Start(ctx, "image::GetImageBytes")
	defer span.End()

	meta, err := s.GetImageMetadata(ctx, id)
	if err != nil {
		return nil, err
	}

	resourceId := meta.OriginalResourceId
	if thumbnail {
		resourceId = meta.ThumbnailResourceId
	}

	assetBytes, err := s.assetStorage.GetAssetBytes(ctx, BucketName, resourceId)
	if err != nil {
		log.Error().Err(err).Msgf("failed getting image asset bytes from storage (id=%s, resourceId=%s, thumbnail=%t)", id, resourceId, thumbnail)
		return nil, ErrInternal
	}

	log.Info().Msgf("read image asset (id=%s, resourceId=%s, thumbnail=%t)", id, resourceId, thumbnail)
	if thumbnail {
		s.thumbDownloadsCounter.Inc()
	} else {
		s.origDownloadsCounter.Inc()
	}

	return assetBytes, nil
}
