package assets

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

func newFileRepo(s3 *minio.Client, tracer trace.Tracer) *fileRepo {
	return &fileRepo{
		s3:     s3,
		tracer: tracer,
	}
}

type fileRepo struct {
	s3     *minio.Client
	tracer trace.Tracer
}

func (f *fileRepo) SaveFile(ctx context.Context, bucket, id string, rBytes []byte) error {
	_, span := f.tracer.Start(ctx, "s3::SaveFile")
	defer span.End()

	opts := minio.PutObjectOptions{} //ContentType: "image/jpeg"}
	_, err := f.s3.PutObject(ctx, bucket, id, bytes.NewReader(rBytes), int64(len(rBytes)), opts)
	if err != nil {
		return err
	}

	return nil
}

func (f *fileRepo) GetFile(ctx context.Context, bucket, id string) ([]byte, error) {
	_, span := f.tracer.Start(ctx, "s3::GetFile")
	defer span.End()

	//TODO: return err only when db access fails
	obj, err := f.s3.GetObject(ctx, bucket, id, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
