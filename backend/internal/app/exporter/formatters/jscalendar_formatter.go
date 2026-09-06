package formatters

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/potibm/tidsapparat/internal/app/domain"
)

type JsCalendarFormatter struct {
	ProductID     string
	Timezone      string
	DefaultAdress string
}

func NewJsCalendarFormatter(productID, timezone, defaultAddress string) *JsCalendarFormatter {
	if productID == "" {
		productID = "-//Tidsapparat//Timetable//EN"
	}

	if timezone == "" {
		timezone = "Europe/Berlin"
	}

	return &JsCalendarFormatter{ProductID: productID, Timezone: timezone, DefaultAdress: defaultAddress}
}

func (f *JsCalendarFormatter) Extension() string {
	return ".json"
}

func (f *JsCalendarFormatter) ContentType() string {
	return "application/json; charset=utf-8"
}

func (f *JsCalendarFormatter) Format(entries domain.TimeTable) ([]byte, error) {
	var latestUpdate time.Time
	for _, entry := range entries {
		if entry.UpdatedAt.After(latestUpdate) {
			latestUpdate = entry.UpdatedAt
		}
	}

	group := JSGroup{
		Type:    "Group",
		UID:     "tidsapparat-timetable-group",
		ProdID:  f.ProductID,
		Updated: latestUpdate.UTC().Format(time.RFC3339),
		Entries: make([]JSEvent, 0, len(entries)),
	}

	for _, entry := range entries {
		event := f.mapToJSEvent(*entry)
		group.Entries = append(group.Entries, event)
	}

	return json.MarshalIndent(group, "", "  ")
}

func (f *JsCalendarFormatter) mapToJSEvent(entry domain.ScheduleEntry) JSEvent {
	uid := fmt.Sprintf("event-%d@demoparty.org", entry.ID)

	var sequence uint64
	if entry.UpdatedAt.After(entry.CreatedAt) {
		sequence = 1
	}

	event := JSEvent{
		Type:        "Event",
		UID:         uid,
		Created:     entry.CreatedAt.UTC().Format(time.RFC3339),
		Updated:     entry.UpdatedAt.UTC().Format(time.RFC3339),
		Sequence:    sequence,
		Title:       entry.Title,
		Description: entry.Description,
		Start:       entry.StartTime.Format("2006-01-02T15:04:05"),
		TimeZone:    f.Timezone,
		Duration:    formatDuration(entry.EndTime.Sub(entry.StartTime)),
		Status:      "confirmed",
	}

	if entry.ExternalURL != "" {
		event.Links = map[string]JSLink{
			"link-1": {
				Type: "Link",
				Href: entry.ExternalURL,
			},
		}
	}

	if entry.Location != nil {
		address := f.DefaultAdress

		if entry.Location.Address != nil && *entry.Location.Address != "" {
			address = *entry.Location.Address
		}

		loc := JSLocation{
			Type: "Location",
			Name: entry.Location.Name,
		}

		if address != "" {
			loc.Description = address
		}

		event.Locations = map[string]JSLocation{
			"loc-1": loc,
		}
	}

	if entry.Category != nil {
		event.Keywords = map[string]bool{
			entry.Category.Name: true,
		}
	}

	return event
}

func formatDuration(d time.Duration) string {
	secs := int(d.Seconds())
	if secs <= 0 {
		return "PT0S"
	}

	const (
		secondsPerMinute = 60
		secondsPerHour   = 3600
	)

	h := secs / secondsPerHour
	m := (secs % secondsPerHour) / secondsPerMinute
	s := secs % secondsPerMinute

	res := "PT"
	if h > 0 {
		res += fmt.Sprintf("%dH", h)
	}

	if m > 0 {
		res += fmt.Sprintf("%dM", m)
	}

	if s > 0 {
		res += fmt.Sprintf("%dS", s)
	}

	return res
}
