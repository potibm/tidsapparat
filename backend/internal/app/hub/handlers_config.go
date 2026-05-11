package hub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/potibm/tidsapparat/internal/app/config"
)

type AppConfigPublic struct {
	Version            string                         `json:"version"`
	Environment        string                         `json:"environment"`
	EnvironmentMessage string                         `json:"environment_message"`
	DateLocale         string                         `json:"date_locale"`
	DateOptions        config.DateFormatOptionsConfig `json:"date_options"`
	Sentry             config.SentryConfig            `json:"sentry"`
	Timezone           string                         `json:"timezone"`
	PartyDays          []PartyDaysPublic              `json:"party_days"`
	EventDurations     []int                          `json:"event_durations"`
}

type PartyDaysPublic struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Server) handleGetPublicConfig(c *gin.Context) {
	pub := mapToPublicConfig(&s.cfg)

	c.JSON(http.StatusOK, pub)
}

func mapToPublicConfig(cfg *config.Config) AppConfigPublic {
	return AppConfigPublic{
		Version:            cfg.App.Version,
		Environment:        cfg.App.Environment,
		EnvironmentMessage: cfg.App.EnvironmentMessage,
		DateLocale:         cfg.Format.Date.Locale,
		DateOptions:        cfg.Format.Date.Options,
		Sentry:             cfg.Sentry,
		Timezone:           cfg.Party.Timezone,
		PartyDays:          generatePartyDaysOrEmpty(cfg.Party.StartDate, cfg.Party.EndDate),
		EventDurations:     cfg.EventDurations,
	}
}

func generatePartyDaysOrEmpty(startDate, endDate string) []PartyDaysPublic {
	partyDays, err := GeneratePartyDays(startDate, endDate)
	if err != nil {
		return []PartyDaysPublic{}
	}

	return partyDays
}

func GeneratePartyDays(startDateStr, endDateStr string) ([]PartyDaysPublic, error) {
	const (
		hoursInDay                    = 24
		numberOfDaysWhenToAddMoreInfo = 7
		layout                        = "2006-01-02"
	)

	start, err := time.Parse(layout, startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %v", err)
	}

	end, err := time.Parse(layout, endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %v", err)
	}

	if start.After(end) {
		return []PartyDaysPublic{}, nil
	}

	daysDiff := int(end.Sub(start).Hours()/hoursInDay) + 1
	partyDays := make([]PartyDaysPublic, 0, daysDiff)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		name := d.Weekday().String()

		if daysDiff > numberOfDaysWhenToAddMoreInfo {
			name = fmt.Sprintf("%s, %s", name, formatOrdinal(d))
		}

		partyDays = append(partyDays, PartyDaysPublic{
			ID:   d.Format(layout),
			Name: name,
		})
	}

	return partyDays, nil
}

func formatOrdinal(t time.Time) string {
	day := t.Day()
	suffix := "th"

	//nolint:mnd // These are calendar days for ordinal formatting
	switch day {
	case 1, 21, 31:
		suffix = "st"
	case 2, 22:
		suffix = "nd"
	case 3, 23:
		suffix = "rd"
	}

	return fmt.Sprintf("%d%s %s", day, suffix, t.Month().String())
}
