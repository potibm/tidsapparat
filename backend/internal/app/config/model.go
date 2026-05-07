//nolint:lll // struct tags can get long and it's more readable to keep them in one line
package config

type RedisURL string

type SentryConfig struct {
	DSN                     string  `json:"dsn"                        mapstructure:"dsn"                        validate:"omitempty,url"`
	TraceSampleRate         float64 `json:"trace_sample_rate"          mapstructure:"trace_sample_rate"          validate:"omitempty,gte=0,lte=1"`
	ReplaySessionSampleRate float64 `json:"replay_session_sample_rate" mapstructure:"replay_session_sample_rate" validate:"omitempty,gte=0,lte=1"`
	ReplayErrorSampleRate   float64 `json:"replay_error_sample_rate"   mapstructure:"replay_error_sample_rate"   validate:"omitempty,gte=0,lte=1"`
	Environment             string  `json:"environment"                mapstructure:"environment"                validate:"required"`
	Version                 string  `json:"version"                    mapstructure:"version"                    validate:"required"`
}

type AppConfig struct {
	Version string `mapstructure:"version"`

	GinMode     string `mapstructure:"gin_mode" validate:"required,oneof=debug release test"`
	Environment string `mapstructure:"env"      validate:"required,oneof=development staging production test"`

	LogLevel  string `mapstructure:"log_level"  validate:"required,oneof=debug info warn error"`
	LogFormat string `mapstructure:"log_format" validate:"required,oneof=json text"`

	DbFilename         string                 `mapstructure:"db_filename"         validate:"required"`
	FrontendURL        string                 `mapstructure:"frontend_url"        validate:"required,http_url"`
	CorsAllowOrigins   CorsAllowOriginsConfig `mapstructure:"cors_allow_origins"  validate:"dive,required"`
	EnvironmentMessage string                 `mapstructure:"environment_message"`
	RedisURL           RedisURL               `mapstructure:"redis_url"           validate:"omitempty,url"`
}

type CorsAllowOriginsConfig []string

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Sentry SentryConfig `mapstructure:"sentry"`
}
