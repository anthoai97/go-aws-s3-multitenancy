package storage_s3

import (
	"context"
	"time"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *StorageS3) CreateFolder(ctx context.Context, folderPath string, client *s3.Client) (*string, error) {
	// Timeout hanlde
	ctx, cancel := context.WithTimeout(ctx, core.REQUEST_TIMEOUT*time.Second)
	defer cancel()

	respcn := make(chan core.ResponseChan[*string])

	go func() {
		input := &s3.PutObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(folderPath),
		}

		_, err := client.PutObject(ctx, input)

		if err != nil {
			respcn <- core.ResponseChan[*string]{
				Data:  nil,
				Error: err,
			}
			return
		}

		respcn <- core.ResponseChan[*string]{
			Data:  aws.String("Create folder successfully"),
			Error: err,
		}
	}()

	// Case: return forward
	for {
		select {
		case <-ctx.Done():
			return nil, core.ErrTimeout

		case resp := <-respcn:
			return resp.Data, resp.Error
		}
	}

}
