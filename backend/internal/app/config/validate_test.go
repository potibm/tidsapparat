package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_PlaylistDefaultsAndValidation(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			GinMode:     "debug",
			Environment: "development",
			LogLevel:    "info",
			LogFormat:   "text",
			DbFilename:  "test.db",
			FrontendURL: "http://localhost:3000",
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
	}

	// 1. trigger validation
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestAppConfig_Validate(t *testing.T) {
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

	cfg.DbFilename = "../invalid-filename"
	err = cfg.Validate()
	assert.Error(t, err)
}
