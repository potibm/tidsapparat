package initializer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/potibm/billedapparat/internal/app/config"
)

func InitializeS3Client(
	ctx context.Context,
	s3ClientConfig *config.S3ClientConfig,
	logger *slog.Logger,
) (*s3.Client, error) {
	if s3ClientConfig == nil {
		return nil, nil
	}

	awsCfg, err := SetupAWSConfig(
		ctx,
		s3ClientConfig.AccessKeyID,
		s3ClientConfig.SecretAccessKey,
		s3ClientConfig.Region,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s3ClientConfig.Endpoint)
		o.UsePathStyle = s3ClientConfig.UsePathStyle
	})

	return s3Client, nil
}
