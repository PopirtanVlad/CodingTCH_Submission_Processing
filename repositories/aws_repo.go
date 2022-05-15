package repositories

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"os"
	"strings"
)

var (
	s3session *s3.S3
)

const (
	BUCKET_NAME = "vladbucket123"
	REGION      = "eu-central-1"
)

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))
}

func uploadObject(fileName string) (resp *s3.PutObjectOutput) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Uploading: ", fileName)
	resp, err = s3session.PutObject(&s3.PutObjectInput{Body: f,
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(strings.Split(fileName, "/")[1]),
		ACL:    aws.String(s3.BucketCannedACLPublicRead)})

	if err != nil {
		panic(err)
	}
	return resp
}

func getObject(filename string) {
	fmt.Println("Downloading:", filename)

	resp, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
	})

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		panic(err)
	}
}

func deleteObject(filename string) (resp *s3.DeleteObjectOutput) {
	fmt.Println("Deleting", filename)

	resp, err := s3session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
	})

	if err != nil {
		panic(err)
	}

	return resp
}
