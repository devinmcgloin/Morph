package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/devinmcgloin/morph/src/env"
)

func main() {
	region := env.Getenv("AWS-REGION", "us-east-1")
	bucket := env.Getenv("AWS_S3_BUCKET", "morph-content")

	svc := s3.New(session.New(&aws.Config{Region: aws.String(region)}))

	log.Printf("Creating bucket %s in region %s", bucket, region)
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: &bucket,
	})
	if err != nil {
		log.Fatal("Failed to create bucket", err)
		return
	}

	if err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &bucket}); err != nil {
		log.Printf("Failed to wait for bucket to exist %s, %s\n", bucket, err)
		return
	}

}
