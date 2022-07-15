package core

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client

func initAws() {
	S3Client = s3.New(s3.Options{
		Region: Settings.AWS_DEFAULT_REGION,
		Credentials: aws.CredentialsProviderFunc(func(c context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     Settings.AWS_ACCESS_KEY_ID,
				SecretAccessKey: Settings.AWS_SECRET_ACCESS_KEY,
			}, nil
		}),
	})
}
