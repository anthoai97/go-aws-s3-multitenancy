package token_vendor_machine

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/ethereum/go-ethereum/log"
)

type TemplateValue struct {
	Bucket string
	Tenant string
}

type TokenVendorMachine struct {
	RoleArn             string
	BucketName          string
	ExternalID          string
	CredDurationSeconds int32
	Logger              log.Logger
}

func NewTokenVendorMachine(RoleArn, ExternalID, BucketName string, CredDurationSeconds int32, Logger log.Logger) *TokenVendorMachine {
	return &TokenVendorMachine{
		RoleArn:             RoleArn,
		ExternalID:          ExternalID,
		Logger:              Logger,
		BucketName:          BucketName,
		CredDurationSeconds: CredDurationSeconds,
	}
}

func (mc *TokenVendorMachine) genegrateVendorPolicy(tenantID, bucket, templatePath string) (string, error) {
	mc.Logger.Debug("Generate new Vendor Policy", "bucket", bucket, "tenantID", tenantID, "templatePath", templatePath)
	tmpl := template.Must(template.ParseFiles(templatePath))

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, TemplateValue{Bucket: bucket, Tenant: tenantID})
	if err != nil {
		mc.Logger.Debug("Generate new Vendor policy error", "tenantID", tenantID, "err", err)
		return "", err
	}
	return buf.String(), nil
}

func (mc *TokenVendorMachine) RequestVendorCredentials(ctx context.Context, tenantID string) (*types.Credentials, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	svc := sts.NewFromConfig(cfg)

	policy, err := mc.genegrateVendorPolicy(tenantID, mc.BucketName, "./templates/access-s3-customer-storage.json")
	if err != nil {
		return nil, err
	}
	mc.Logger.Debug("RequestVendorCredentials", "genegrateVendorPolicy", policy)
	fmt.Println(policy)
	roleSessionName := fmt.Sprintf("%s-%d", tenantID, time.Now().Unix())
	assumeRoleInput := sts.AssumeRoleInput{
		RoleArn:         aws.String(mc.RoleArn),
		Policy:          aws.String(policy),
		RoleSessionName: aws.String(roleSessionName),
		DurationSeconds: aws.Int32(mc.CredDurationSeconds), // 900s minimum requrired value
		ExternalId:      aws.String(mc.ExternalID),
	}
	result, err := svc.AssumeRole(ctx, &assumeRoleInput)
	if err != nil {
		return nil, fmt.Errorf("failed to assume role, %v", err)
	}

	return result.Credentials, nil
}
