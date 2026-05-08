package writers

import (
	"bytes"
	"context"
	"fmt"
	"strings"

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

func (w *S3Writer) Write(ctx context.Context, filename string, data []byte) error {
	body := bytes.NewReader(data)

	_, err := w.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(w.bucket),
		Key:         aws.String(filename),
		Body:        body,
		ContentType: aws.String(w.determineContentType(filename)),
	})
	if err != nil {
		return fmt.Errorf("s3 upload failed: %w", err)
	}

	return nil
}

func (w *S3Writer) determineContentType(filename string) string {
	if strings.HasSuffix(filename, ".ics") {
		return "text/calendar"
	}

	if strings.HasSuffix(filename, ".json") {
		return "application/json"
	}

	return "application/octet-stream"
}
