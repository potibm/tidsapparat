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
	CorsAllowOrigins   CorsAllowOriginsConfig `mapstructure:"cors_allow_origins"  validate:"required,dive,required"`
	EnvironmentMessage string                 `mapstructure:"environment_message"`
	RedisURL           RedisURL               `mapstructure:"redis_url"           validate:"omitempty,url"`
}

type ExporterConfig struct {
	Name        string            `mapstructure:"name"        validate:"required"`
	Type        string            `mapstructure:"type"        validate:"required,oneof=ical"`
	Destination string            `mapstructure:"destination" validate:"required,oneof=s3 file"`
	Filename    string            `mapstructure:"filename"    validate:"required"`
	Options     map[string]string `mapstructure:"options"`
	Enabled     bool              `mapstructure:"enabled"`
}

type CorsAllowOriginsConfig []string

type PartyConfig struct {
	Timezone       string `mapstructure:"timezone"        validate:"required,timezone"`
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

type FormatConfig struct {
	Date DateFormatConfig `json:"date" mapstructure:"date"`
}

type DateFormatOptionsConfig map[string]any

type DateFormatConfig struct {
	Locale  string                  `json:"locale"  mapstructure:"locale"  validate:"required"`
	Options DateFormatOptionsConfig `json:"options" mapstructure:"options"`
}

type Config struct {
	App            AppConfig        `mapstructure:"app"`
	Format         FormatConfig     `mapstructure:"format"`
	Sentry         SentryConfig     `mapstructure:"sentry"`
	Exporter       []ExporterConfig `mapstructure:"exporter"`
	Party          PartyConfig      `mapstructure:"party"`
	S3Client       *S3ClientConfig  `mapstructure:"s3_client"`
	EventDurations []int            `mapstructure:"event_durations" validate:"dive,gte=0"`
}
