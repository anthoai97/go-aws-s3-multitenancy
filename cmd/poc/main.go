package main

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Tenant struct {
	Tenant string
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		msg := fmt.Sprintf("failed to create new session, %v", err)
		fmt.Println(msg)
	}

	svc := sts.New(sess)
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
		DurationSeconds: aws.Int64(60 * 60 * 1),
		ExternalId:      aws.String(externalId),
	}

	result, err := svc.AssumeRole(&assumeRoleInput)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role, %v", err)
		fmt.Println(msg)
	}

	s3Session, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(
			*result.Credentials.AccessKeyId,
			*result.Credentials.SecretAccessKey,
			*result.Credentials.SessionToken),
	})

	fmt.Println(s3Session)

	s3Service := s3.New(s3Session)
	req, output := s3Service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String("dsr-customer-storage-dev"),
		Key:    aws.String("dataspire/hello.png"), //will be generic via lambda call tenant id parameter
	})
	if err := req.Send(); err != nil {
		fmt.Println(err)
		return
	}
	body, _ := io.ReadAll(output.Body)
	fmt.Println(body)
}
