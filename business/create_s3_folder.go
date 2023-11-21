package business

import (
	"context"
	"strings"

	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (biz *business) CreateS3Folder(ctx context.Context, req *entity.RequestCreateFolder) (*string, error) {
	cred, err := biz.LoadSTSCredentialClaims(ctx.(*gin.Context))
	if err != nil {
		return nil, err
	}

	client, err := biz.S3.GenerateS3Client(ctx, cred)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(req.Path, "/") {
		req.Path = "/" + req.Path
	}

	if !strings.HasSuffix(req.Path, "/") {
		req.Path = req.Path + "/"
	}

	biz.Logger.Debug("CreateS3Folder", "path", req.Tenant+req.Path)
	resp, err := biz.S3.CreateFolder(ctx, req.Tenant+req.Path, client)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
