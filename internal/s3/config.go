package s3

import (
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Config struct {
	Host                 string `env:"HOST"`
	Port                 int    `env:"PORT"`
	AccessKey            string `env:"ACCESS_KEY"`
	SecretAccessKey      string `env:"SECRET_ACCESS_KEY"`
	UseSSL               bool   `env:"USE_SSL"`
	MaxRetries           int    `env:"MAX_RETRIES"`
	PresignedUrlDuration string `env:"PRESIGNED_URL_DURATION"`
}

func (cfg *S3Config) GetPresignedUrlDuration() (time.Duration, error) {
	duration, err := time.ParseDuration(cfg.PresignedUrlDuration)
	if err != nil {
		return time.Duration(0), err
	}

	return duration, nil
}

func (cfg *S3Config) GetOptions() *minio.Options {
	return &minio.Options{
		Creds:      credentials.NewStaticV4(cfg.AccessKey, cfg.SecretAccessKey, ""),
		Secure:     cfg.UseSSL,
		MaxRetries: cfg.MaxRetries,
	}
}

func (cfg *S3Config) GetEndpoint() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
