package business

import (
	"context"

	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
)

func (biz *business) DeleteS3Objects(ctx context.Context, paths []*entity.RequestObjectDelete) (*string, error) {
	return nil, nil
}
