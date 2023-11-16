package business

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (biz *business) DeleteS3Folder(ctx context.Context, req *entity.RequestFolderDelete) (*string, error) {
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

	objectPaths, err := biz.getAllObjects(ctx, req.Tenant+req.Path, client)
	if err != nil {
		return nil, err
	}

	biz.Logger.Debug("DeleteS3Folder", "objectPaths len", len(objectPaths))
	biz.Logger.Debug("DeleteS3Folder", "objectPaths", objectPaths)

	if len(objectPaths) < 1 {
		return nil, fmt.Errorf("can not find folder")
	}

	return biz.S3.DeleteObjects(ctx, objectPaths, client)
}
