package storage_s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *StorageS3) DownloadObjectByGenerateUrl(ctx context.Context, path string, client *s3.PresignClient) (*string, error) {
	req, err := client.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(3600 * int64(time.Second)) // 1 Hour
	})

	if err != nil {
		return nil, err
	}

	return &req.URL, nil
}
