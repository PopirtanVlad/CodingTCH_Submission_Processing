package repositories

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

const (
	BUCKET_NAME = "vladbucket123"
	REGION      = "eu-central-1"
)

type S3Repository struct {
	s3session *s3.S3
}

func NewS3Repository() *S3Repository {
	return &S3Repository{
		s3session: s3.New(session.Must(session.NewSession(&aws.Config{Region: aws.String(REGION)}))),
	}
}

func (s3Repo *S3Repository) uploadObject(fileName string) (resp *s3.PutObjectOutput) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	logrus.WithFields(logrus.Fields{
		"File Name": fileName,
	}).Info("Uploading file to s3")
	resp, err = s3Repo.s3session.PutObject(&s3.PutObjectInput{Body: f,
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(strings.Split(fileName, "/")[1]),
		ACL:    aws.String(s3.BucketCannedACLPublicRead)})

	if err != nil {
		panic(err)
	}
	return resp
}

func (s3Repo *S3Repository) GetSubmission(fileName string) (io.ReadCloser, error) {
	logrus.WithFields(logrus.Fields{
		"File Name": fileName,
	}).Info("Downloading submission from s3")

	resp, err := s3Repo.s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (s3Repo *S3Repository) deleteObject(fileName string) (resp *s3.DeleteObjectOutput) {
	logrus.WithFields(logrus.Fields{
		"File Name": fileName,
	}).Info("Deleting file from s3")

	resp, err := s3Repo.s3session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(fileName),
	})

	if err != nil {
		panic(err)
	}

	return resp
}
