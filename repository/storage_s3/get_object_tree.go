package storage_s3

import (
	"context"
	"time"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (s *StorageS3) GetObjectTree(ctx context.Context, path string, client *s3.Client) (data *entity.S3ObjectTree, err error) {
	// Timeout hanlde
	ctx, cancel := context.WithTimeout(ctx, core.REQUEST_TIMEOUT*time.Second)
	defer cancel()

	path = core.FormatBucketPrefixForTree(path)
	respcn := make(chan core.ResponseChan[*entity.S3ObjectTree])

	s.Logger.Debug("ListS3StorageTree", "Path", path)
	s.Logger.Debug("ListS3StorageTree", "Bucket", s.Bucket)

	// Excute
	go func() {
		nextToken, ok := ctx.Value("NextContinuationToken").(*string)
		if !ok {
			nextToken = nil
		}

		params := &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.Bucket),
			Prefix:            aws.String(path),
			Delimiter:         aws.String("/"),
			MaxKeys:           21,
			ContinuationToken: nextToken,
		}

		req, err := client.ListObjectsV2(ctx, params)

		if err != nil {
			respcn <- core.ResponseChan[*entity.S3ObjectTree]{
				Data:  nil,
				Error: err,
			}
			return
		}

		tree := &entity.S3ObjectTree{}

		// Load data
		for _, pref := range req.CommonPrefixes {
			// element is the element from someSlice for where we are
			tree.CommonPrefixes = append(tree.CommonPrefixes, pref.Prefix)
		}

		for _, cont := range req.Contents {
			tree.Contents = append(tree.Contents, &entity.S3Object{
				Key:          cont.Key,
				Size:         cont.Size,
				LastModified: cont.LastModified,
			})
		}

		tree.Prefix = req.Prefix

		tree.NextContinuationToken = req.NextContinuationToken

		respcn <- core.ResponseChan[*entity.S3ObjectTree]{
			Data:  tree,
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
