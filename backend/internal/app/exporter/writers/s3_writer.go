package writers

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Writer struct {
	client *s3.Client
	bucket string
}

func NewS3Writer(client *s3.Client, bucket string) *S3Writer {
	return &S3Writer{
		client: client,
		bucket: bucket,
	}
}

func (w *S3Writer) Write(ctx context.Context, filename string, data []byte, contentType string) error {
	body := bytes.NewReader(data)

	_, err := w.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(w.bucket),
		Key:         aws.String(filename),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("s3 upload failed: %w", err)
	}

	return nil
}
