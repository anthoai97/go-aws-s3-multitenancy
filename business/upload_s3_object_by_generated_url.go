package business

import (
	"context"
	"strings"
	"sync"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/ptr"
	"github.com/gin-gonic/gin"
)

func (biz *business) UploadS3ObjectsByGenerateUrl(ctx context.Context, objects []*entity.RequestFileUpload) ([]*entity.ResponseFileUpload, error) {
	cred, err := biz.LoadSTSCredentialClaims(ctx.(*gin.Context))
	if err != nil {
		return nil, err
	}

	client, err := biz.S3.GeneratePresignClient(ctx, cred)
	if err != nil {
		return nil, err
	}

	var resp = []*entity.ResponseFileUpload{}

	// validate and process
	for _, file := range objects {
		if len(file.Path) < 1 || len(file.Tenant) < 1 {
			return nil, core.ErrBadRequest
		}

		if !strings.HasPrefix(file.Path, "/") {
			file.Path = "/" + file.Path
		}
	}

	resChan := make(chan *string, 1000)
	var wg = new(sync.WaitGroup)

	for _, file := range objects {
		wg.Add(1)
		go func(file *entity.RequestFileUpload, s3client *s3.PresignClient) {
			defer wg.Done()

			biz.Logger.Debug("UploadS3ObjectsByGenerateUrl", "path", file.Tenant+file.Path)
			url, err := biz.S3.UploadObjectByGenerateUrl(ctx, file.Tenant+file.Path, biz.defaultObjectTagging(file.Tenant), s3client)
			if err != nil {
				resChan <- ptr.String(err.Error())
				return
			}

			resChan <- url

		}(file, client)
	}

	for _, file := range objects {
		url := <-resChan
		resp = append(resp, &entity.ResponseFileUpload{
			Tenant:         file.Tenant,
			Path:           file.Path,
			UploadUrl:      url,
			HeaderTagKey:   "x-amz-tagging",
			HeaderTagValue: biz.defaultObjectTagging(file.Tenant),
		})
	}

	// close the channel since every data is sent.
	wg.Wait()
	close(resChan)

	return resp, nil
}
