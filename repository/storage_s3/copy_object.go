package storage_s3

import (
	"context"
	"fmt"
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
		s.Logger.Debug("CopyObject", "path", fmt.Sprintf("%v/%v", s.Bucket, path), "newPath", newPath)
		input := &s3.CopyObjectInput{
			Bucket:           aws.String(s.Bucket),
			CopySource:       aws.String(fmt.Sprintf("%v/%v", s.Bucket, path)),
			Key:              aws.String(newPath),
			TaggingDirective: types.TaggingDirectiveCopy,
		}

		_, err := client.CopyObject(ctx, input)

		if err != nil {
			s.Logger.Debug("CopyObject", "err", err)
			respcn <- core.ResponseChan[bool]{
				Data:  false,
				Error: err,
			}
			return
		}

		s.Logger.Debug("CopyObject", "data", "Copy object successful")
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
