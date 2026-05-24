package config

import (
	"strings"
	"time"

	"github.com/potibm/tidsapparat/internal/app/calendar"
	"github.com/spf13/viper"
)

const (
	OtelServiceName        = "tidsapparat"
	OtelBackendServiceName = OtelServiceName + "-backend"

	DefaultPort = 8080

	DefaultTraceSampleRate         = 0.1
	DefaultReplaySessionSampleRate = 0.1
	DefaultReplayErrorSampleRate   = 0.1

	DataDirname = "./data"

	DefaultDBFilename = "tidsapparat"

	DataDirPerm = 0o755
)

var DefaultDateOptions = DateFormatOptionsConfig{
	"weekday": "short",
	"hour":    "2-digit",
	"minute":  "2-digit",
}

func InitViper() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("app.port", DefaultPort)
	viper.SetDefault("app.otel_endpoint", "")
	viper.SetDefault("app.gin_mode", "release")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.env", "production")
	viper.SetDefault("app.db_filename", DefaultDBFilename)
	viper.SetDefault("app.frontend_url", "")
	viper.SetDefault("app.cors_allow_origins", []string{})
	viper.SetDefault("app.redis_url", "")

	viper.SetDefault("format.date.locale", "da-DK")
	viper.SetDefault("format.date.options", DefaultDateOptions)

	viper.SetDefault("sentry.dsn", "")
	viper.SetDefault("sentry.trace_sample_rate", DefaultTraceSampleRate)
	viper.SetDefault("sentry.replay_session_sample_rate", DefaultReplaySessionSampleRate)
	viper.SetDefault("sentry.replay_error_sample_rate", DefaultReplayErrorSampleRate)

	viper.SetDefault("party.timezone", "Europe/Copenhagen")
	viper.SetDefault("party.start_date", calendar.GetWeekdayCurrentWeek(time.Friday).Format("2006-01-02"))
	viper.SetDefault("party.end_date", calendar.GetWeekdayCurrentWeek(time.Sunday).Format("2006-01-02"))

	viper.SetDefault("event_durations", []int{0, 15, 30, 60, 90, 120})

	viper.RegisterAlias("sentry.environment", "app.env")
	viper.RegisterAlias("sentry.version", "app.version")
}
