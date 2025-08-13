package s3

import (
	"chat_app_backend/internal/extensions"
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type Statuses = string

const (
	ObjectNotFound Statuses = "NoSuchKey"
	BucketNotFound          = "NoSuchBucket"
)

type Buckets = string

const (
	AvatarsBucket       Buckets = "avatars"
	InterestsIconBucket         = "interests"
)

type FileType = string

const (
	Png  FileType = "png"
	Jpeg          = "jpeg"
)

type IClient interface {
	GetDownloadUrl(ctx context.Context, filename string, bucketName Buckets) (string, error)
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, filename string, bucketName Buckets) (string, error)
	DeleteFile(ctx context.Context, filename string, bucketName Buckets) error
	ModifyFileContents(ctx context.Context, fileHeader *multipart.FileHeader, filename string, bucketName Buckets) (string, error)
	CreateBucket(ctx context.Context, bucketName Buckets) error
	BucketExists(ctx context.Context, bucketName Buckets) (bool, error)
	FileExists(ctx context.Context, filename string, bucketName Buckets) (bool, error)
}

type Client struct {
	client               *minio.Client
	presignedUrlDuration time.Duration
}

func (c *Client) GetDownloadUrl(ctx context.Context, filename string, bucketName Buckets) (string, error) {
	exists, err := c.FileExists(ctx, filename, bucketName)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("file %s does not exist in bucket %s", filename, bucketName)
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	urlObject, urlGettingError := c.client.PresignedGetObject(ctx, bucketName, filename, c.presignedUrlDuration, reqParams)
	if urlGettingError != nil {
		return "", urlGettingError
	}

	return urlObject.String(), nil
}

func (c *Client) FileExists(ctx context.Context, filename string, bucketName Buckets) (bool, error) {
	_, err := c.client.StatObject(ctx, bucketName, filename, minio.StatObjectOptions{})
	if err == nil {
		return true, nil
	}

	errResponse := minio.ToErrorResponse(err)
	if errResponse.Code != ObjectNotFound {
		return false, err
	}

	return false, nil
}

func (c *Client) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, filename string, bucketName Buckets) (string, error) {
	exists, err := c.FileExists(ctx, filename, bucketName)
	if err != nil {
		return "", err
	}

	if exists {
		return "", fmt.Errorf("file %s already exists in bucket %s", filename, bucketName)
	}

	file, openFileError := fileHeader.Open()
	if openFileError != nil {
		return "", openFileError
	}

	defer func(f multipart.File) {
		_ = f.Close()
	}(file)

	_, uploadError := c.client.PutObject(
		ctx,
		bucketName,
		filename,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType:     fileHeader.Header.Get("Content-Type"),
			ContentEncoding: fileHeader.Header.Get("Content-Transfer-Encoding"),
		},
	)

	if uploadError != nil {
		return "", uploadError
	}

	return c.GetDownloadUrl(ctx, filename, bucketName)
}

func (c *Client) DeleteFile(ctx context.Context, filename string, bucketName Buckets) error {
	exists, err := c.FileExists(ctx, filename, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("file %s does not exist in bucket %s", filename, bucketName)
	}

	return c.client.RemoveObject(ctx, bucketName, filename, minio.RemoveObjectOptions{})
}

func (c *Client) ModifyFileContents(ctx context.Context, fileHeader *multipart.FileHeader, filename string, bucketName Buckets) (string, error) {
	exists, err := c.FileExists(ctx, filename, bucketName)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("file %s does not exist in bucket %s", filename, bucketName)
	}

	file, fileOpeningError := fileHeader.Open()
	if fileOpeningError != nil {
		return "", fileOpeningError
	}

	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	_, uploadError := c.client.PutObject(
		ctx,
		bucketName,
		filename,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{
			ContentType:     fileHeader.Header.Get("Content-Type"),
			ContentEncoding: fileHeader.Header.Get("Content-Transfer-Encoding"),
		},
	)

	if uploadError != nil {
		return "", uploadError
	}

	return c.GetDownloadUrl(ctx, filename, bucketName)
}

func (c *Client) CreateBucket(ctx context.Context, bucketName Buckets) error {
	exists, getExistenceError := c.client.BucketExists(ctx, bucketName)
	if getExistenceError != nil {
		return getExistenceError
	}

	if exists {
		return fmt.Errorf("bucket with name %s already exists", bucketName)
	}

	return c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}

func (c *Client) BucketExists(ctx context.Context, bucketName Buckets) (bool, error) {
	return c.client.BucketExists(ctx, bucketName)
}

func CreateClient(cfg *S3Config) (*Client, error) {
	client, err := minio.New(cfg.GetEndpoint(), cfg.GetOptions())
	if err != nil {
		return nil, err
	}

	duration, durationParseError := cfg.GetPresignedUrlDuration()
	if durationParseError != nil {
		return nil, durationParseError
	}

	return &Client{
		client:               client,
		presignedUrlDuration: duration,
	}, nil
}

func ConstructFilenameFromFileType(fileType FileType) string {
	return fmt.Sprintf("%s.%s", extensions.NewUUID(), fileType)
}
