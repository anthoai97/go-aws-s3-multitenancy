package business

import (
	"context"
	"strings"

	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (biz *business) DeleteS3Objects(ctx context.Context, paths []*entity.RequestObjectDelete) (*string, error) {
	cred, err := biz.LoadSTSCredentialClaims(ctx.(*gin.Context))
	if err != nil {
		return nil, err
	}

	client, err := biz.S3.GenerateS3Client(ctx, cred)
	if err != nil {
		return nil, err
	}

	var objectPaths []string
	err = nil
	for _, obj := range paths {
		path := obj.Path

		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		biz.Logger.Debug("DeleteS3Objects", "path", obj.Tenant+path)
		objectPaths = append(objectPaths, obj.Tenant+path)
	}

	if err != nil {
		return nil, err
	}

	return biz.S3.DeleteObjects(ctx, objectPaths, client)
}
