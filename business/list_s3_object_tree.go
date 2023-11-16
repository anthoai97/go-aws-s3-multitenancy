package business

import (
	"context"
	"strings"

	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (biz *business) ListS3StorageTree(ctx context.Context, path, tenant string) (*entity.S3ObjectTree, error) {
	cred, err := biz.LoadSTSCredentialClaims(ctx.(*gin.Context))
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	client, err := biz.S3.GenerateS3Client(ctx, cred)
	if err != nil {
		return nil, err
	}

	res, err := biz.S3.GetObjectTree(ctx, tenant+path, client)
	if err != nil {
		return nil, err
	}

	return res, nil
}
