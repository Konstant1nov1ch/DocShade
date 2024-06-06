package s3_provider

import (
	"bytes"
	"context"
	"errors"
	"time"

	"gitlab.com/docshade/common/core"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	Bucket                     = "preprocessing"
	minioFiLeNotFoundErrorCode = "NoSuchKey"
	maxRetries                 = 10
	retryDelay                 = 5 * time.Second
)

type S3 interface {
	// InitS3 инициализировать s3
	InitS3() error
	// Put загружает файл в S3
	Put(ctx context.Context, objectName, path string, objectBody []byte, metaData map[string]string) error
	// IsObjectExist проверяет, существует ли объект в S3
	IsObjectExist(ctx context.Context, path, objectName string) (bool, error)
	// CreateBucket создает бакет в S3, если он не существует
	CreateBucket(ctx context.Context, bucketName string) error
}

type s3 struct {
	cfg core.S3Config
	s3  *minio.Client
}

func NewS3(cfg core.S3Config) S3 {
	return &s3{
		cfg: cfg,
	}
}

func (s *s3) InitS3() error {
	var err error

	for i := 0; i < maxRetries; i++ {
		s.s3, err = minio.New(
			s.cfg.Endpoint,
			&minio.Options{
				Creds:  credentials.NewStaticV4(s.cfg.AccessKeyID, s.cfg.SecretAccessKey, ""),
				Secure: false, // Использование HTTPS
			})

		if err == nil {
			break
		}

		time.Sleep(retryDelay)
	}

	return err
}

func (s *s3) Put(ctx context.Context, objectName, path string, objectBody []byte, metaData map[string]string) error {
	err := s.resolvePath(ctx, path)
	if err != nil {
		return err
	}
	isObjectExist, err := s.IsObjectExist(ctx, path, objectName)
	if err != nil {
		return err
	}
	if isObjectExist {
		return errors.New("File with name '" + objectName + "' in bucket '" + path + "' already exists")
	}

	buffer := bytes.NewBuffer(objectBody)
	objLength := int64(len(objectBody))
	_, err = s.s3.PutObject(
		ctx,
		path,
		objectName,
		buffer,
		objLength,
		minio.PutObjectOptions{UserMetadata: metaData},
	)
	return err
}

func (s *s3) resolvePath(ctx context.Context, path string) error {
	return s.CreateBucket(ctx, path)
}

func (s *s3) IsObjectExist(ctx context.Context, path, objectName string) (bool, error) {
	_, err := s.s3.StatObject(
		ctx,
		path,
		objectName,
		minio.StatObjectOptions{},
	)

	if err != nil {
		errorRes := minio.ToErrorResponse(err)
		if errorRes.Code == minioFiLeNotFoundErrorCode {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *s3) CreateBucket(ctx context.Context, bucketName string) error {
	isBucketExist, err := s.s3.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	if !isBucketExist {
		return s.s3.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}
	return nil
}
