package storage_s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ethereum/go-ethereum/log"
)

type StorageS3 struct {
	// SVC           *s3.Client
	Bucket string
	Logger log.Logger
}

func NewStorageS3(bucket string, logger log.Logger) *StorageS3 {
	return &StorageS3{
		Bucket: bucket,
		Logger: logger,
	}
}

func (*StorageS3) generateS3Client(ctx context.Context, cred *aws.CredentialsCache) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	cfg.Credentials = cred
	client := s3.NewFromConfig(cfg)

	return client, nil
}

func (storage *StorageS3) generatePresignClient(ctx context.Context, cred *aws.CredentialsCache) (*s3.PresignClient, error) {
	s3Client, err := storage.generateS3Client(ctx, cred)
	if err != nil {
		return nil, err
	}

	return s3.NewPresignClient(s3Client), nil
}
