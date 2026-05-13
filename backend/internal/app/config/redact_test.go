package config

import (
	"testing"
)

func TestRedisURL_Redacted(t *testing.T) {
	const expectedRedacted = "%2A%2A%2AREDACTED%2A%2A%2A"

	tests := []struct {
		name string
		url  RedisURL
		want RedisURL
	}{
		{
			name: "URL with password should be redacted",
			url:  "redis://user:password123@localhost:6379",
			want: "redis://user:" + expectedRedacted + "@localhost:6379",
		},
		{
			name: "URL with only username should be redacted",
			url:  "redis://admin@localhost:6379",
			want: "redis://admin:" + expectedRedacted + "@localhost:6379",
		},
		{
			name: "URL without user info should remain unchanged",
			url:  "redis://localhost:6379",
			want: "redis://localhost:6379",
		},
		{
			name: "Invalid URL should be returned as is",
			url:  "://invalid-url",
			want: "://invalid-url",
		},
		{
			name: "Empty URL should return empty string",
			url:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.url.Redacted(); got != tt.want {
				t.Errorf("RedisURL.Redacted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func assertS3ClientRedacted(t *testing.T, got, want *S3ClientConfig) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Errorf("S3Client = %v, want nil", got)
		}

		return
	}

	if got == nil {
		t.Fatalf("S3Client is nil, want non-nil")
	}

	if got.AccessKeyID != want.AccessKeyID {
		t.Errorf("AccessKeyID = %v, want %v", got.AccessKeyID, want.AccessKeyID)
	}

	if got.SecretAccessKey != want.SecretAccessKey {
		t.Errorf(
			"SecretAccessKey = %v, want %v",
			got.SecretAccessKey,
			want.SecretAccessKey,
		)
	}
}

func TestConfig_RedactConfigForDisplay(t *testing.T) {
	tests := []struct {
		name         string
		config       Config
		wantS3Client *S3ClientConfig
		wantRedisURL RedisURL
	}{
		{
			name: "nil S3Client does not panic and fields are redacted",
			config: Config{
				Sentry: SentryConfig{DSN: "secret-dsn"},
				App:    AppConfig{RedisURL: "redis://user:pass@localhost:6379"},
			},
			wantS3Client: nil,
			wantRedisURL: "redis://user:%2A%2A%2AREDACTED%2A%2A%2A@localhost:6379",
		},
		{
			name: "non-nil S3Client gets redacted",
			config: Config{
				S3Client: &S3ClientConfig{
					AccessKeyID:     "access-key",
					SecretAccessKey: "secret-key",
				},
				App: AppConfig{RedisURL: "redis://localhost:6379"},
			},
			wantS3Client: &S3ClientConfig{
				AccessKeyID:     "***REDACTED***",
				SecretAccessKey: "***REDACTED***",
			},
			wantRedisURL: "redis://localhost:6379",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.RedactConfigForDisplay()

			assertS3ClientRedacted(t, result.S3Client, tt.wantS3Client)

			if result.App.RedisURL != tt.wantRedisURL {
				t.Errorf("RedisURL = %v, want %v", result.App.RedisURL, tt.wantRedisURL)
			}
		})
	}
}
