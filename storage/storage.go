package storage

import (
	"context"
	"git.kuschku.de/justjanne/imghost-frontend/configuration"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/url"
	"os"
	"time"
)

type Storage struct {
	config   configuration.StorageConfiguration
	s3client *minio.Client
}

func NewStorage(config configuration.StorageConfiguration) (storage Storage, err error) {
	storage.config = config
	storage.s3client, err = minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.Secure,
	})
	return
}

func (storage Storage) UploadFile(ctx context.Context, bucketName string, fileName string, mimeType string, file *os.File) (err error) {
	_, err = storage.s3client.FPutObject(
		ctx,
		bucketName,
		fileName,
		file.Name(),
		minio.PutObjectOptions{
			ContentType: mimeType,
		})
	return
}

func (storage Storage) Upload(ctx context.Context, bucketName string, fileName string, mimeType string, reader io.Reader) (err error) {
	_, err = storage.s3client.PutObject(
		ctx,
		bucketName,
		fileName,
		reader,
		-1,
		minio.PutObjectOptions{
			ContentType: mimeType,
		})
	return
}

func (storage Storage) DownloadFile(ctx context.Context, bucketName string, fileName string, file *os.File) (err error) {
	err = storage.s3client.FGetObject(
		ctx,
		bucketName,
		fileName,
		file.Name(),
		minio.GetObjectOptions{})
	return
}

func (storage Storage) UrlFor(ctx context.Context, bucketName string, fileName string) (url *url.URL, err error) {
	url, err = storage.s3client.PresignedGetObject(
		ctx,
		bucketName,
		fileName,
		7*24*time.Hour,
		map[string][]string{})
	return
}
