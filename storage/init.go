package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// creates a session with credentials 
func s3Configure(credential map[string]string) (*session.Session, error) {
	s3Configuration := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(credential["hostkey"], credential["secretkey"], credential["token"]),
		Endpoint:         aws.String(credential["endpoint"]),
		Region:           aws.String(credential["region"]),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Configuration)
	if err != nil {
		return nil, err
	}

	return newSession, nil
}
