package business

import (
	"context"

	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
)

func (biz *business) UploadS3ObjectsByGenerateUrl(ctx context.Context, objects []*entity.RequestFileUpload) ([]*entity.ResponseFileUpload, error) {
	return nil, nil
}
