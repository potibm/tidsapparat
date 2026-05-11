package config

import (
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	validDbFilename = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	validLocale     = regexp.MustCompile(`^[a-zA-Z]{2}-[A-Z]{2}$`)
)

func (c *Config) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	if err := c.App.Validate(); err != nil {
		return err
	}

	return nil
}

func (f *AppConfig) Validate() error {
	if !validDbFilename.MatchString(f.DbFilename) {
		return fmt.Errorf("db_filename '%s' contains invalid characters", f.DbFilename)
	}

	return nil
}

func (ru *RedisURL) Validate() error {
	rString := string(*ru)

	if !ru.IsValid() {
		return fmt.Errorf("redis_url '%s' is not a valid URL", rString)
	}

	redisURL := ru.URLObject()
	if redisURL.Scheme != "redis" && redisURL.Scheme != "rediss" {
		return fmt.Errorf(
			"redis_url '%s' has invalid scheme '%s' (expected 'redis' or 'rediss')",
			rString,
			redisURL.Scheme,
		)
	}

	host, _, err := net.SplitHostPort(redisURL.Host)
	if err != nil || host == "" {
		return fmt.Errorf("redis_url '%s' has missing host", rString)
	}

	return nil
}

func (c *PartyConfig) Validate() error {
	const (
		maxPartyDuration = 31
		hoursInDay       = 24
	)

	layout := "2006-01-02"

	start, err := time.Parse(layout, c.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date format")
	}

	end, err := time.Parse(layout, c.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end_date format")
	}

	if start.After(end) {
		return fmt.Errorf("start_date (%s) must be before or equal to end_date (%s)", c.StartDate, c.EndDate)
	}

	if end.Before(time.Now().AddDate(-1, 0, 0)) {
		return fmt.Errorf("event date is more than a year in the past")
	}

	days := int(end.Sub(start).Hours()/hoursInDay) + 1
	if days > maxPartyDuration {
		return fmt.Errorf("party duration exceeds maximum limit of %d days (current: %d)", maxPartyDuration, days)
	}

	return nil
}
