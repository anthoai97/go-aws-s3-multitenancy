package storage_s3

import (
	"context"
	"time"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *StorageS3) ListObject(ctx context.Context, path, nextContinuationToken string, client *s3.Client) (data *s3.ListObjectsV2Output, err error) {
	// Timeout hanlde
	ctx, cancel := context.WithTimeout(ctx, core.REQUEST_TIMEOUT*time.Second)
	defer cancel()

	respcn := make(chan core.ResponseChan[*s3.ListObjectsV2Output])

	// Excute
	go func() {
		var nextToken *string

		if len(nextContinuationToken) > 0 {
			nextToken = aws.String(nextContinuationToken)
		}

		params := &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.Bucket),
			Prefix:            aws.String(path),
			MaxKeys:           1000,
			ContinuationToken: nextToken,
		}

		req, err := client.ListObjectsV2(ctx, params)

		if err != nil {
			respcn <- core.ResponseChan[*s3.ListObjectsV2Output]{
				Data:  nil,
				Error: err,
			}
			return
		}

		respcn <- core.ResponseChan[*s3.ListObjectsV2Output]{
			Data:  req,
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
