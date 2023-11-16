package business

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

func (biz *business) GenerateSTSCredential(ctx context.Context, tenant string) (*types.Credentials, error) {
	cred, err := biz.TVM.RequestVendorCredentials(ctx, tenant)
	if err != nil {
		biz.Logger.Debug("GenerateSTSCredential", "Error", err)
		return nil, err
	}

	return cred, nil
}
