package storage

import (
	"context"
	"creatly-task/internal/config"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Storage struct {
	connection *s3.S3
	timeout    time.Duration
	bucketName string
}

type Config struct {
	Timeout    time.Duration
	Region     string
	BucketName string
	AccessKey  string
	SecretKey  string
}

func New(cfg *config.Storage) (*Storage, error) {

	creds := credentials.NewStaticCredentialsFromCreds(credentials.Value{
		AccessKeyID:     cfg.AccessKey,
		SecretAccessKey: cfg.SecretKey,
	})

	config := aws.NewConfig().WithCredentials(creds).WithRegion(cfg.Region)
	session, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	svc := s3.New(session)

	_, err = svc.ListObjects(&s3.ListObjectsInput{
		Bucket: &cfg.BucketName,
	})
	if err != nil {
		return nil, err
	}

	return &Storage{
		connection: svc,
		timeout:    cfg.Timeout,
		bucketName: cfg.BucketName,
	}, nil
}

func (s *Storage) UploadFile(file multipart.File, bucket, key string) error {
	ctx, _ := context.WithTimeout(context.Background(), s.timeout)
	out, err := s.connection.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	fmt.Println(out)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
		}
	}
	return err
}
