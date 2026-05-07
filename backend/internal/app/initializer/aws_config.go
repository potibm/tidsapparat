package initializer

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func SetupAWSConfig(ctx context.Context, key, secret, region string) (aws.Config, error) {
	staticCreds := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(key, secret, ""),
	)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(staticCreds),
	)

	return cfg, err
}
