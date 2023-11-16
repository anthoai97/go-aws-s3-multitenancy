package storage_s3

import (
	"context"
	"time"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (s *StorageS3) DeleteObjects(ctx context.Context, paths []string, client *s3.Client) (*string, error) {
	// Timeout hanlde
	ctx, cancel := context.WithTimeout(ctx, core.REQUEST_TIMEOUT*time.Second)
	defer cancel()

	respcn := make(chan core.ResponseChan[*string])

	go func() {
		var objectIds []types.ObjectIdentifier
		for _, key := range paths {
			s.Logger.Debug("DeleteObjects", "key", key)
			objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
		}

		input := &s3.DeleteObjectsInput{
			Bucket: aws.String(s.Bucket),
			Delete: &types.Delete{Objects: objectIds},
		}

		_, err := client.DeleteObjects(ctx, input)

		if err != nil {
			respcn <- core.ResponseChan[*string]{
				Data:  nil,
				Error: err,
			}
			return
		}

		respcn <- core.ResponseChan[*string]{
			Data:  aws.String("Deleted Successfully"),
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
