package AWS

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadImageAWS(file []byte, size int64, filename string, bucketURI string, region string) (string, error) {

	svc := s3.New(session.New(&aws.Config{Region: aws.String(region)}))
	destPath := "/content/" + filename
	params, err := formatParams(file, size, bucketURI, destPath)

	if err != nil {
		log.Printf("Error while creating AWS params %s", err)
		return "", err
	}

	_, err = svc.PutObject(params)
	if err != nil {
		log.Printf("Error while uploading to aws %s", err)
		return "", err
	}

	return fmt.Sprintf("https://s3.amazonaws.com/%s%s", bucketURI, destPath), nil
}

func formatParams(buffer []byte, size int64, bucketName string, path string) (*s3.PutObjectInput, error) {

	fileBytes := bytes.NewReader(buffer)

	fileType := http.DetectContentType(buffer)

	log.Printf("Uploading %s to %s with size %d and type %s", path, bucketName, size, fileType)

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	return params, nil

}
