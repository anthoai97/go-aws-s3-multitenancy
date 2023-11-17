package business

import (
	"context"
	"strings"

	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (biz *business) DownloadS3ObjectsByGenerateUrl(ctx context.Context, object *entity.RequestFileDownload) (*entity.ResponseFilDownload, error) {
	cred, err := biz.LoadSTSCredentialClaims(ctx.(*gin.Context))
	if err != nil {
		return nil, err
	}

	client, err := biz.S3.GeneratePresignClient(ctx, cred)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(object.Path, "/") {
		object.Path = "/" + object.Path
	}

	url, err := biz.S3.DownloadObjectByGenerateUrl(ctx, object.Tenant+object.Path, client)
	if err != nil {
		return nil, err
	}

	return &entity.ResponseFilDownload{
		Url:    url,
		Tenant: object.Tenant,
		Path:   object.Path,
	}, nil
}
