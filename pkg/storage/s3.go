package storage

import (
	"bytes"
	"creatly-task/internal/config"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Storage struct {
	connection *s3.S3
	timeout    time.Duration
	bucketName string
	region     string
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
		region:     cfg.Region,
	}, nil
}

func (s *Storage) UploadFile(file []byte, filesize int64, filename string) (string, error) {

	_, err := s.connection.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s.bucketName),
		Key:                  aws.String(filename),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(file),
		ContentLength:        aws.Int64(filesize),
		ContentType:          aws.String(http.DetectContentType(file)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	fileUrl := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", s.bucketName, s.region, filename)
	return fileUrl, err
}
