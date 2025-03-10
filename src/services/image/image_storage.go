package image

import (
	"bytes"
	"context"

	"github.com/minio/minio-go/v7"
)

func NewStorage() *fileStorage {
	return &fileStorage{}
}

type fileStorage struct {
	s3 *minio.Client
}

func (s *fileStorage) Save(ctx context.Context, name string, img []byte) error {
	bucketName := "test"

	opts := minio.PutObjectOptions{ContentType: "image/jpeg"}
	_, err := s.s3.PutObject(ctx, bucketName, name, bytes.NewReader(img), int64(len(img)), opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *fileStorage) Get(ctx context.Context, name string) ([]byte, error) {
	bucketName := "test"

	//TODO: return err only when db access fails
	obj, err := s.s3.GetObject(ctx, bucketName, name, minio.GetObjectOptions{})
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
