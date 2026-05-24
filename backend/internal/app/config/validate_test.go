package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_PlaylistDefaultsAndValidation(t *testing.T) {
	currentYear := time.Now().Format("2006")

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
			StartDate: currentYear + "-01-01",
			EndDate:   currentYear + "-01-02",
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

func TestPartyConfig_Validate(t *testing.T) {
	currentYear := time.Now().Format("2006")

	tests := []struct {
		name    string
		config  PartyConfig
		wantErr bool
	}{
		{
			name: "Valid range - same day",
			config: PartyConfig{
				StartDate: currentYear + "-05-01",
				EndDate:   currentYear + "-05-01",
			},
			wantErr: false,
		},
		{
			name: "Valid range - multiple days",
			config: PartyConfig{
				StartDate: currentYear + "-05-01",
				EndDate:   currentYear + "-05-05",
			},
			wantErr: false,
		},
		{
			name: "Invalid range - far too long ago",
			config: PartyConfig{
				StartDate: "1997-05-01",
				EndDate:   "1997-05-03",
			},
			wantErr: true,
		},
		{
			name: "Invalid - end before start",
			config: PartyConfig{
				StartDate: currentYear + "-05-10",
				EndDate:   currentYear + "-05-01",
			},
			wantErr: true,
		},
		{
			name: "Invalid - wrong date format",
			config: PartyConfig{
				StartDate: "01-05-2026",
				EndDate:   currentYear + "-05-10",
			},
			wantErr: true,
		},
		{
			name: "Invalid - way too long (over 31 days)",
			config: PartyConfig{
				StartDate: currentYear + "-05-01",
				EndDate:   currentYear + "-07-01",
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
