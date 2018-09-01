package storage

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/png"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fokal/fokal-core/pkg/log"
)

type AWSStorageService struct {
	bucketURI string
	region    string
	kind      string // content, avatar
}

func New(bucket, region, kind string) *AWSStorageService {
	return &AWSStorageService{bucketURI: bucket, region: region, kind: kind}
}

func (ss *AWSStorageService) UploadImage(ctx context.Context, img image.Image, shortcode string) error {
	var err error

	path := strings.Join([]string{ss.kind, shortcode}, "/")
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return err
	}

	err = ss.imageAWS(ctx, buf, "png", path)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return errors.New("Error while uploading image")
	}
	return nil
}

func (ss *AWSStorageService) DeleteImage(ctx context.Context, shortcode string) error {
	return errors.New("Not Implemented")
}

func (ss *AWSStorageService) imageAWS(ctx context.Context, img *bytes.Buffer, format string, filename string) error {

	sess, err := session.NewSession(&aws.Config{Region: aws.String(ss.region)})
	if err != nil {
		log.WithContext(ctx).Errorf("error while constructing new aws session %s", err)
		return err
	}
	svc := s3.New(sess)

	params, err := ss.formatParams(ctx, img, int64(img.Len()), format, filename)

	if err != nil {
		log.WithContext(ctx).Errorf("Error while creating AWS params %s", err)
		return err
	}

	_, err = svc.PutObject(params)
	if err != nil {
		log.WithContext(ctx).Errorf("Error while uploading to aws %s", err)
		return err
	}

	return nil
}

func (ss *AWSStorageService) formatParams(ctx context.Context, buffer *bytes.Buffer, size int64, filetype string, path string) (*s3.PutObjectInput, error) {

	fileBytes := bytes.NewReader(buffer.Bytes())

	log.WithContext(ctx).Infof("Uploading %s to %s with size %d and type %s", path, ss.bucketURI, size, filetype)

	params := &s3.PutObjectInput{
		Bucket:        aws.String(ss.bucketURI),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String("image/" + filetype),
	}

	return params, nil
}
