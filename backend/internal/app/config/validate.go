package config

import (
	"fmt"
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

	if err := c.Format.Validate(); err != nil {
		return err
	}

	if err := c.Party.Validate(); err != nil {
		return err
	}

	return nil
}

func (f *AppConfig) Validate() error {
	if !validDbFilename.MatchString(f.DbFilename) {
		return fmt.Errorf("db_filename '%s' contains invalid characters", f.DbFilename)
	}

	if f.RedisURL != "" {
		if err := f.RedisURL.Validate(); err != nil {
			return err
		}
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

func (f *FormatConfig) Validate() error {
	return f.Date.Validate()
}

func (f *DateFormatConfig) Validate() error {
	if !validLocale.MatchString(f.Locale) {
		return fmt.Errorf("date.locale '%s' is not a valid locale", f.Locale)
	}

	return nil
}
