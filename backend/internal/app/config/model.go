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

type ExporterConfig struct {
	Name        string            `mapstructure:"name"`
	Type        string            `mapstructure:"type"`
	Destination string            `mapstructure:"destination"`
	Filename    string            `mapstructure:"filename"`
	Options     map[string]string `mapstructure:"options"`
}

type CorsAllowOriginsConfig []string

type PartyConfig struct {
	Timezone       string `mapstructure:"timezone"        validate:"required"`
	DefaultAddress string `mapstructure:"default_address"`
	StartDate      string `mapstructure:"start_date"      validate:"required,datetime=2006-01-02"`
	EndDate        string `mapstructure:"end_date"        validate:"required,datetime=2006-01-02"`
}

type S3ClientConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id"     validate:"required"`
	SecretAccessKey string `mapstructure:"secret_access_key" validate:"required"`
	Region          string `mapstructure:"region"            validate:"required"`
	Endpoint        string `mapstructure:"endpoint"          validate:"required"`
	UsePathStyle    bool   `mapstructure:"use_path_style"`
}

type Config struct {
	App      AppConfig        `mapstructure:"app"`
	Sentry   SentryConfig     `mapstructure:"sentry"`
	Exporter []ExporterConfig `mapstructure:"exporter"`
	Party    PartyConfig      `mapstructure:"party"`
	S3Client *S3ClientConfig  `mapstructure:"s3_client"`
}
