package uploadservice

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/png"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fokal/fokal-core/pkg/logger"
)

var mediaTypeOptions = []string{"jp2", "jpeg", "png", "tiff", "bmp"}

type AWSStorageService struct {
	bucketURI string
	region    string
	kind      string // content, avatar
}

func (ss AWSStorageService) UploadImage(ctx context.Context, img image.Image, shortcode string) error {
	var err error

	path := strings.Join([]string{ss.kind, shortcode}, "/")
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}

	err = imageAWS(buf, "png", path, "images-fokal", "us-west-1")
	if err != nil {
		logger.Error(ctx, err)
		err := errors.New("Error while uploading image")
		return err
	}
	return nil
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
