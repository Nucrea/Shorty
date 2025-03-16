package s3

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel/trace"
)

func NewFileStorage(s3 *minio.Client, tracer trace.Tracer, bucket string) *FileStorage {
	return &FileStorage{
		bucket: bucket,
		s3:     s3,
		tracer: tracer,
	}
}

type FileStorage struct {
	bucket string
	s3     *minio.Client
	tracer trace.Tracer
}

func (f *FileStorage) SaveFile(ctx context.Context, name string, img []byte) error {
	_, span := f.tracer.Start(ctx, "s3::SaveFile")
	defer span.End()

	opts := minio.PutObjectOptions{} //ContentType: "image/jpeg"}
	_, err := f.s3.PutObject(ctx, f.bucket, name, bytes.NewReader(img), int64(len(img)), opts)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStorage) GetFile(ctx context.Context, name string) ([]byte, error) {
	_, span := f.tracer.Start(ctx, "s3::GetFile")
	defer span.End()

	//TODO: return err only when db access fails
	obj, err := f.s3.GetObject(ctx, f.bucket, name, minio.GetObjectOptions{})
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
