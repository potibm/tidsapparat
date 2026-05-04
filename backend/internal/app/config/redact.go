package config

const redacted = "***REDACTED***"

func (c Config) RedactConfigForDisplay() Config {
	result := c

	result.Sentry.DSN = redacted

	return result
}
