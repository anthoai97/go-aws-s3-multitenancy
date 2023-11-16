package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Tenant struct {
	Tenant string
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		msg := fmt.Sprintf("failed to create new session, %v", err)
		fmt.Println(msg)
	}

	svc := sts.NewFromConfig(cfg)

	tmpl := template.Must(template.ParseFiles("./template.json"))
	tenantId := "dataspire"        // will be parametric
	externalId := "someexternalid" // will be taken from config file

	var buf bytes.Buffer
	tmpl.Execute(&buf, Tenant{Tenant: tenantId})

	policy := buf.String()
	fmt.Println(policy)

	roleArn := "arn:aws:iam::360307311296:role/ri.developer-assume-role" // will be taken from config file
	roleSessionName := fmt.Sprintf("%s-%d", tenantId, time.Now().Unix())
	assumeRoleInput := sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		Policy:          aws.String(policy),
		RoleSessionName: aws.String(roleSessionName),
		DurationSeconds: aws.Int32(60 * 60 * 1),
		ExternalId:      aws.String(externalId),
	}

	result, err := svc.AssumeRole(context.TODO(), &assumeRoleInput)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role, %v", err)
		fmt.Println(msg)
	}

	// s3Session, _ := session.NewSession(&aws.Config{
	// 	Region: aws.String("ap-southeast-1"),
	// 	Credentials: credentials.NewStaticCredentials(
	// 		*result.Credentials.AccessKeyId,
	// 		*result.Credentials.SecretAccessKey,
	// 		*result.Credentials.SessionToken),
	// })

	cfg2, _ := config.LoadDefaultConfig(context.TODO())

	cred := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(*result.Credentials.AccessKeyId,
		*result.Credentials.SecretAccessKey,
		*result.Credentials.SessionToken))

	cfg2.Credentials = cred
	client := s3.NewFromConfig(cfg2)

	req, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("dsr-customer-storage-dev"),
		Prefix: aws.String("dataspire/"),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	// body, _ := io.ReadAll(req.Body)
	fmt.Println(req.Contents)
}
