package upload

import (
	"bytes"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadImageAWS(img *bytes.Buffer, format string, filename string, bucketURI string, region string) error {

	svc := s3.New(session.New(&aws.Config{Region: aws.String(region)}))

	params, err := formatParams(img, int64(img.Len()), format, bucketURI, filename)

	if err != nil {
		log.Printf("Error while creating AWS params %s", err)
		return err
	}

	_, err = svc.PutObject(params)
	if err != nil {
		log.Printf("Error while uploading to aws %s", err)
		return err
	}

	return nil
}

func formatParams(buffer *bytes.Buffer, size int64, filetype string, bucketName string, path string) (*s3.PutObjectInput, error) {

	fileBytes := bytes.NewReader(buffer.Bytes())

	log.Printf("Uploading %s to %s with size %d and type %s", path, bucketName, size, filetype)

	params := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String("image/" + filetype),
	}

	return params, nil
}
