package config

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	OtelServiceName        = "tidsapparat"
	OtelBackendServiceName = OtelServiceName + "-backend"

	DefaultTraceSampleRate         = 0.1
	DefaultReplaySessionSampleRate = 0.1
	DefaultReplayErrorSampleRate   = 0.1

	DataDirname = "./data"

	DefaultDBFilename = "tidsapparat"

	DataDirPerm = 0o755
)

func InitViper() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("app.gin_mode", "release")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.env", "production")
	viper.SetDefault("app.db_filename", DefaultDBFilename)
	viper.SetDefault("app.frontend_url", "")
	viper.SetDefault("app.cors_allow_origins", []string{})
	viper.SetDefault("app.redis_url", "")

	viper.SetDefault("sentry.dsn", "")
	viper.SetDefault("sentry.trace_sample_rate", DefaultTraceSampleRate)
	viper.SetDefault("sentry.replay_session_sample_rate", DefaultReplaySessionSampleRate)
	viper.SetDefault("sentry.replay_error_sample_rate", DefaultReplayErrorSampleRate)

	viper.RegisterAlias("sentry.environment", "app.env")
	viper.RegisterAlias("sentry.version", "app.version")
}
