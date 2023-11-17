package storage_s3

import (
	"context"
	"time"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (s *StorageS3) CopyObject(ctx context.Context, path, newPath string, client *s3.Client) (bool, error) {
	// Timeout hanlde
	ctx, cancel := context.WithTimeout(ctx, core.REQUEST_TIMEOUT*time.Second)
	defer cancel()

	respcn := make(chan core.ResponseChan[bool])

	go func() {
		input := &s3.CopyObjectInput{
			Bucket:           aws.String(s.Bucket),
			CopySource:       &path,
			Key:              &newPath,
			TaggingDirective: types.TaggingDirectiveCopy,
		}

		_, err := client.CopyObject(ctx, input)

		if err != nil {
			respcn <- core.ResponseChan[bool]{
				Data:  false,
				Error: err,
			}
			return
		}

		respcn <- core.ResponseChan[bool]{
			Data:  true,
			Error: err,
		}
	}()

	// Case: return forward
	for {
		select {
		case <-ctx.Done():
			return false, core.ErrTimeout

		case resp := <-respcn:
			return resp.Data, resp.Error
		}
	}

}
