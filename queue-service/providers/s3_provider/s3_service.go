package s3_provider

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"gitlab.com/docshade/common/core"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	BucketIn                   = "preprocessing"
	minioFiLeNotFoundErrorCode = "NoSuchKey"
	BucketOut                  = "postprocessing"
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
	// Remove удаляет файл из S3
	Remove(ctx context.Context, objectName, path string) error
	// Move перемещает файл из одного бакета в другой
	Move(ctx context.Context, objectName, srcPath, destPath, newDirName string) (string, error)
	// Get извлекает файл из S3
	Get(ctx context.Context, objectName, path string) ([]byte, error)
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

// InitS3 инициализирует S3 клиент
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

func (s *s3) Remove(ctx context.Context, path string, objectName string) error {
	err := s.s3.RemoveObject(ctx, BucketIn, objectName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return err
	}
	return nil
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
	// log.Println(objectName)
	// log.Println(BucketOut)
	buffer := bytes.NewBuffer(objectBody)
	objLength := int64(len(objectBody))
	_, err = s.s3.PutObject(
		ctx,
		BucketOut,
		objectName+".pdf",
		buffer,
		objLength,
		minio.PutObjectOptions{UserMetadata: metaData},
	)
	return err
}

func (s *s3) resolvePath(ctx context.Context, path string) error {
	isBucketExist, err := s.s3.BucketExists(ctx, path)
	if err != nil {
		return err
	}
	if !isBucketExist {
		return s.s3.MakeBucket(ctx, path, minio.MakeBucketOptions{})
	}
	return nil
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

func (s *s3) Move(ctx context.Context, objectName, srcPath, destPath, newDirName string) (string, error) {
	// Get the object from the source bucket
	object, err := s.s3.GetObject(ctx, srcPath, objectName, minio.GetObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get object: %v", err)
	}
	defer object.Close()

	// Read the object data
	objectData, err := io.ReadAll(object)
	if err != nil {
		return "", fmt.Errorf("failed to read object data: %v", err)
	}

	// Create the new directory in the destination bucket
	fullDestPath := fmt.Sprintf("%s/%s", destPath, newDirName)
	err = s.resolvePath(ctx, fullDestPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %v", err)
	}

	// Put the object into the destination bucket
	err = s.Put(ctx, objectName, fullDestPath, objectData, nil)
	if err != nil {
		return "", fmt.Errorf("failed to put object: %v", err)
	}

	// Remove the object from the source bucket
	err = s.Remove(ctx, objectName, srcPath)
	if err != nil {
		return "", fmt.Errorf("failed to remove object: %v", err)
	}

	return fmt.Sprintf("%s/%s", fullDestPath, objectName), nil
}

func (s *s3) Get(ctx context.Context, objectName, path string) ([]byte, error) {
	// log.Println(path)

	object, err := s.s3.GetObject(ctx, BucketIn, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}
	defer object.Close()

	objectData, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %v", err)
	}

	return objectData, nil
}
