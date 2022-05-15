package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

func listBuckets() (resp *s3.ListBucketsOutput) {

	resp, err := s3session.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		panic(err)
	}

	return resp

}

func createBucket() (resp *s3.CreateBucketOutput) {
	resp, err := s3session.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(BUCKET_NAME),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(REGION),
		},
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println("Bucket name already in use!")
				panic(err)
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println("Bucket already exists in your possession!")
			default:
				panic(err)
			}
		}
	}
	return resp
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

func listObjects() (resp *s3.ListObjectsV2Output) {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(BUCKET_NAME),
	})
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

func uploadFiles() {
	folder := "testfile"

	files, _ := ioutil.ReadDir(folder)

	for _, file := range files {
		if !file.IsDir() {
			uploadObject(folder + "/" + file.Name())
		}
	}
}

func main() {
	uploadObject("testfile/solution1.txt")
	fmt.Println(listObjects())
	getObject("solution1.txt")
	deleteObject("solution1.txt")
	fmt.Println(listObjects())
}
