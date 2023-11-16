package api

import (
	"context"

	"github.com/anthoai97/go-aws-s3-multitenancy/entity"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

type Business interface {
	ListS3StorageTree(ctx context.Context, path, tenant string) (*entity.S3ObjectTree, error)
	RenameS3Object(ctx context.Context, path string) (*entity.S3ObjectTree, error)
	DeleteS3Objects(ctx context.Context, paths []*entity.RequestObjectDelete) (*string, error)
	UploadS3ObjectsByGenerateUrl(ctx context.Context, objects []*entity.RequestFileUpload) ([]*entity.ResponseFileUpload, error)
	DownloadS3ObjectsByGenerateUrl(ctx context.Context, objects []string) (*string, error)
	GenerateSTSCredential(ctx context.Context, tenant string) (*types.Credentials, error)
	DeleteS3Folder(sctx context.Context, req *entity.RequestFolderDelete) (*string, error)
}

type api struct {
	business Business
}

func NewAPI(biz Business) *api {
	return &api{
		business: biz,
	}
}
