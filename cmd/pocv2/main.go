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
	// cred := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("ASIAVHY75ILAD2FJG2NH",
	// 	"1rAm3GE1IzTKdUhCOTR1z+PWo15xJWurnvE5VKdf",
	// 	"IQoJb3JpZ2luX2VjEDwaDmFwLXNvdXRoZWFzdC0xIkcwRQIhAJJSqwJW9S2Fhf+/fF4x3DKJbbmy51gbnsh0EIWPUxptAiAl08/mnpdus/zkqQlkh3M5z9M9U2OmraAkZyfTEZ04UCq3AwiF//////////8BEAMaDDM2MDMwNzMxMTI5NiIMyTPXba0SV0KaH4cxKosDac2y4czo7nGwejyyBS+fGcF2hT2g3n7HeBMIByk4NeM+4lF/Hh1OpP2nPm75UrBhBBWNknRT4S9dXXQTdpwAhOxT151PCFQb/17k1GANZO+OjPNw/F10LhzbU07YqBDxCS4xKuERo7ZYg+Zq7xwZaVDt1uMpRLrbYSZvO5otioBVpM2atoPHfff3Llpsp6q9DvCP4sO0pZZkA/2Jsf/CY5PFOGFqb2x69ombYzSfZUUU4roZOkQrH7W17km/8Ql1sjcxkVWDWyWruuA6FcfjcLa3vU5opALq+B19UJeGFkZOqNgK3XbXSsiVmYI6kgl3Vq0I1h1Xblw7KpzQqvKwS4VJKngCke033hHLLVprHif/V21NtwQoUM1qZvoYlqtgC3F9snBjTZ03e7HzfXGSwTlYp7yGTfLlkgxpu7LQQ5K/97fz5feQLRhKMXM7AxibV+Nw2aJRd6JVL7z2UYrcrnWhb1nbwz6sKSyflpCznYpwNMVePPO8VCNrL5wyvAq+hhr6tU/x7iG+JAwwrabWqgY6nQHzbXWCeDrpnjE7/aSUjx1rgufsW0hLhBtZz43UT3jdOCVvLWruQXGe9EOCryL6InOtF4UTBaSMnOwuvO/+j15rAOqlnFPqyNw0ir8QSGB/g4hN3w0Va6h4IPGWWBMz06efTqhYBNrD93+PYHJtVbV2Sz7b1Du/3taz0jDFM8XTIaK7ekcOAiEw4Q5ImoraCRyveBsCTMl2ITLZ39aI"))

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
