package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_PlaylistDefaultsAndValidation(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			GinMode:          "debug",
			Environment:      "development",
			LogLevel:         "info",
			LogFormat:        "text",
			DbFilename:       "test.db",
			CorsAllowOrigins: []string{"https://localhost:3333", "https://localhost:3121"},
			FrontendURL:      "http://localhost:3000",
		},
		Sentry: SentryConfig{
			DSN:                     "https://test@sentry.io/123",
			TraceSampleRate:         0.1,
			ReplaySessionSampleRate: 0.1,
			ReplayErrorSampleRate:   0.1,
			Environment:             "development",
			Version:                 "1.2.3",
		},
		Party: PartyConfig{
			Timezone:  "Europe/Berlin",
			StartDate: "2024-01-01",
			EndDate:   "2024-01-02",
		},
		Format: FormatConfig{
			Date: DateFormatConfig{
				Locale: "en-US",
				Options: DateFormatOptionsConfig{
					"year":  "numeric",
					"month": "long",
					"day":   "numeric",
				},
			},
		},
		Exporter: []ExporterConfig{},
	}

	// 1. trigger validation
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestAppConfig_Validate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := AppConfig{
			GinMode:     "debug",
			Environment: "development",
			LogLevel:    "info",
			LogFormat:   "text",
			DbFilename:  "test.db",
			FrontendURL: "http://localhost:3000",
		}

		err := cfg.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid db filename", func(t *testing.T) {
		cfg := AppConfig{
			GinMode:     "debug",
			Environment: "development",
			LogLevel:    "info",
			LogFormat:   "text",
			DbFilename:  "../invalid-filename",
			FrontendURL: "http://localhost:3000",
		}

		err := cfg.Validate()
		assert.Error(t, err)
	})

	t.Run("valid redis url", func(t *testing.T) {
		cfg := AppConfig{
			GinMode:     "debug",
			Environment: "development",
			LogLevel:    "info",
			LogFormat:   "text",
			DbFilename:  "test.db",
			FrontendURL: "http://localhost:3000",
			RedisURL:    "redis://localhost:6379",
		}

		err := cfg.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid redis url", func(t *testing.T) {
		cfg := AppConfig{
			GinMode:     "debug",
			Environment: "development",
			LogLevel:    "info",
			LogFormat:   "text",
			DbFilename:  "test.db",
			FrontendURL: "http://localhost:3000",
			RedisURL:    "http://localhost:6379",
		}

		err := cfg.Validate()
		assert.Error(t, err)
	})
}

func TestRedisURL_Validate(t *testing.T) {
	tests := []struct {
		name    string
		url     RedisURL
		wantErr bool
	}{
		{
			name:    "valid redis URL",
			url:     "redis://localhost:6379",
			wantErr: false,
		},
		{
			name:    "valid rediss URL",
			url:     "rediss://localhost:6379",
			wantErr: false,
		},
		{
			name:    "invalid scheme",
			url:     "http://localhost:6379",
			wantErr: true,
		},
		{
			name:    "missing host",
			url:     "redis://:6379",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			url:     "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.url.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPartyConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  PartyConfig
		wantErr bool
	}{
		{
			name: "Valid range - same day",
			config: PartyConfig{
				StartDate: "2026-05-01",
				EndDate:   "2026-05-01",
			},
			wantErr: false,
		},
		{
			name: "Valid range - multiple days",
			config: PartyConfig{
				StartDate: "2026-05-01",
				EndDate:   "2026-05-05",
			},
			wantErr: false,
		},
		{
			name: "Invalid - end before start",
			config: PartyConfig{
				StartDate: "2026-05-10",
				EndDate:   "2026-05-01",
			},
			wantErr: true,
		},
		{
			name: "Invalid - wrong date format",
			config: PartyConfig{
				StartDate: "01-05-2026",
				EndDate:   "2026-05-10",
			},
			wantErr: true,
		},
		{
			name: "Invalid - way too long (over 31 days)",
			config: PartyConfig{
				StartDate: "2026-05-01",
				EndDate:   "2026-07-01",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
