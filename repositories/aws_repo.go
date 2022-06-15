package repositories

import (
	"Licenta_Processing_Service/dtos"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/seqsense/s3sync"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

type S3Repository struct {
	s3session     *s3.S3
	s3Sync        *s3sync.Manager
	baseDirectory string
	bucket        string
}

func NewS3Repository(conf dtos.AWSConfig) *S3Repository {
	s3session := session.Must(session.NewSession(&aws.Config{Region: aws.String(conf.AWSRegion)}))
	return &S3Repository{
		s3session:     s3.New(s3session),
		s3Sync:        s3sync.New(s3session),
		baseDirectory: conf.BaseLocalDir,
		bucket:        conf.AWSBucketName,
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
		Bucket: aws.String(s3Repo.bucket),
		Key:    aws.String(strings.Split(fileName, "/")[1]),
		ACL:    aws.String(s3.BucketCannedACLPublicRead)})

	if err != nil {
		panic(err)
	}
	return resp
}

func (s3Repo *S3Repository) GetSubmission(problemId, submissionId string) (io.ReadCloser, error) {
	logrus.WithFields(logrus.Fields{
		"File Name": submissionId,
	}).Info("Downloading submission from s3")

	filePath := fmt.Sprintf("submissions/%s/%s", problemId, submissionId)
	resp, err := s3Repo.s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Repo.bucket),
		Key:    aws.String(filePath),
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
		Bucket: aws.String(s3Repo.bucket),
		Key:    aws.String(fileName),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func (s3Repo *S3Repository) DownloadTests(problemTitle string) error {
	s3Path := fmt.Sprintf("s3://%s/problems/%s", s3Repo.bucket, problemTitle)
	localPath := fmt.Sprintf("%s/%s", s3Repo.baseDirectory, problemTitle)
	return s3Repo.s3Sync.Sync(s3Path, localPath)
}
