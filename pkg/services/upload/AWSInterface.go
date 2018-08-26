package uploadservice

import (
	"bytes"
	"context"
	"errors"
	"image/jpeg"
	"log"
	"strings"

	"cloud.google.com/go/logging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var mediaTypeOptions = []string{"jp2", "jpeg", "png", "tiff", "bmp"}

type AWSStorageService struct {
	bucketURI string
	region    string
	kind      string // content, avatar
}

func (ss AWSStorageService) UploadImage(ctx context.Context, img *bytes.Buffer, path string) error {
	var err error

	// TODO this does not match properly to the mediaTypeOptions
	in := in(format, mediaTypeOptions)
	if !in {
		err := errors.New("Unsupported Media Type")
		logging.Error(ctx, err)
		return err
	}

	path := strings.Join([]string{kind, shortcode}, "/")
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		logging.Error(ctx, err)
		return err
	}
	err = imageAWS(buf, format, path, "images-fokal", "us-west-1")
	if err != nil {
		logging.Error(ctx, err)
		err := errors.New("Error while uploading image")
		return
	}
	errChan <- nil
}

func (ss AWSStorageService) DeleteImage(ctx context.Context, shortcode string) error {
	return errors.New("Not Implemented")
}

func in(contentType string, opts []string) bool {
	for _, opt := range opts {
		if strings.Compare(contentType, opt) == 0 {
			return true
		}
	}
	return false
}

func imageAWS(img *bytes.Buffer, format string, filename string, bucketURI string, region string) error {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Printf("error while constructing new aws session %s", err)
		return err
	}
	svc := s3.New(sess)

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
