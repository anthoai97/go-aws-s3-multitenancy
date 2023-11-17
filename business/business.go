package business

import (
	"context"

	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	"github.com/anthoai97/go-aws-s3-multitenancy/repository/storage_s3"
	"github.com/anthoai97/go-aws-s3-multitenancy/repository/token_vendor_machine"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"

	md "github.com/anthoai97/go-aws-s3-multitenancy/middleware"
	logger "github.com/ethereum/go-ethereum/log"
)

type business struct {
	TVM    *token_vendor_machine.TokenVendorMachine
	Logger logger.Logger
	S3     *storage_s3.StorageS3
}

func NewBusiness(TVM *token_vendor_machine.TokenVendorMachine, storage_s3 *storage_s3.StorageS3, logger logger.Logger) *business {
	return &business{
		TVM:    TVM,
		Logger: logger,
		S3:     storage_s3,
	}
}

func (*business) defaultObjectTagging(tenant string) string {
	env := core.GetEnvVar("ENV", "develop")
	return "dsrCustomer=" + tenant + "&" + "environment=" + env
}

func (biz *business) LoadSTSCredentialClaims(ctx *gin.Context) (*aws.CredentialsCache, error) {
	access := ctx.GetString(md.CREDENTIAL_ACCESS_KEY)
	secret := ctx.GetString(md.CREDENTIAL_SECRET_KEY)
	session := ctx.GetString(md.CREDENTIAL_SESSION_KEY)
	biz.Logger.Debug("LoadSTSCredentialClaims", "access", access)
	biz.Logger.Debug("LoadSTSCredentialClaims", "secret", secret)
	biz.Logger.Debug("LoadSTSCredentialClaims", "session", session)

	if len(access) < 1 || len(secret) < 1 || len(session) < 1 {
		return nil, core.ErrClaimSTSCredential
	}

	cred := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(access, secret, session))
	biz.Logger.Debug("LoadSTSCredentialClaims", "Step 2", "cred, ok")
	return cred, nil
}

func (biz *business) getAllObjects(ctx context.Context, path string, client *s3.Client) ([]string, error) {
	var next string
	var err error
	var objectPaths []string

	for {
		resp, errp := biz.S3.ListObject(ctx, path, next, client)
		if errp != nil {
			err = errp
			break
		}

		for _, file := range resp.Contents {
			objectPaths = append(objectPaths, *file.Key)
		}

		if resp.IsTruncated {
			next = *resp.NextContinuationToken
		} else {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	return objectPaths, nil
}
