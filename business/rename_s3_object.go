package business

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/gin-gonic/gin"
)

func (biz *business) RenameS3Object(ctx context.Context, object *entity.RequestFileRename) (*entity.ResponseFileRename, error) {
	cred, err := biz.LoadSTSCredentialClaims(ctx.(*gin.Context))
	if err != nil {
		return nil, err
	}

	client, err := biz.S3.GenerateS3Client(ctx, cred)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(object.Path, "/") {
		object.Path = "/" + object.Path
	}

	if !strings.HasPrefix(object.NewPath, "/") {
		object.NewPath = "/" + object.NewPath
	}

	success, err := biz.S3.CopyObject(ctx, object.Tenant+object.Path, object.Tenant+object.NewPath, client)
	if err != nil {
		return nil, err
	}

	// Delete
	if !success {
		return nil, fmt.Errorf("rename object failed")
	}

	_, err = biz.S3.DeleteObjects(ctx, []string{object.Tenant + object.Path}, client)
	if err != nil {
		return nil, err
	}
	// Copy and delete

	return &entity.ResponseFileRename{
		Path:    object.Path,
		NewPath: object.NewPath,
		Tenant:  object.Tenant,
	}, nil
}
