package storage_s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *StorageS3) UploadObjectByGenerateUrl(ctx context.Context, path, tagging string, client *s3.PresignClient) (*string, error) {
	s.Logger.Debug("UploadObjectByGenerateUrl", "tagging", tagging)

	req, err := client.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:  aws.String(s.Bucket),
		Key:     aws.String(path),
		Tagging: aws.String(tagging),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(3600 * int64(time.Second))
	})

	if err != nil {
		return nil, err
	}

	return &req.URL, nil
}
