package formatters

import (
	"strings"
	"testing"
	"time"

	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewIcalFormatter_Defaults(t *testing.T) {
	f := NewIcalFormatter("", "", "")

	assert.Equal(t, "-//Tidsapparat//Timetable//EN", f.ProductID)
	assert.Equal(t, "Europe/Berlin", f.Timezone)
	assert.Empty(t, f.DefaultAdress)
}

func TestNewIcalFormatter_CustomValues(t *testing.T) {
	f := NewIcalFormatter("PRODID", "UTC", "Default Address")

	assert.Equal(t, "PRODID", f.ProductID)
	assert.Equal(t, "UTC", f.Timezone)
	assert.Equal(t, "Default Address", f.DefaultAdress)
}

func TestIcalFormatter_Extension(t *testing.T) {
	f := NewIcalFormatter("", "", "")

	assert.Equal(t, ".ics", f.Extension())
}

func TestIcalFormatter_Format(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)
	address := "123 Main St"

	entries := domain.TimeTable{
		{
			ID:          1,
			Title:       "Opening Ceremony",
			Description: "The grand opening",
			StartTime:   now,
			EndTime:     now.Add(2 * time.Hour),
			CreatedAt:   now,
			UpdatedAt:   now,
			Category: &domain.Category{
				Name:  "Ceremony",
				Color: "#FF0000",
			},
			Location: &domain.Location{
				Name:    "Main Hall",
				Address: &address,
			},
		},
	}

	f := NewIcalFormatter("-//Test//EN", "UTC", "Default Venue")
	result, err := f.Format(entries)
	assert.NoError(t, err)

	output := string(result)
	assert.Contains(t, output, "BEGIN:VCALENDAR")
	assert.Contains(t, output, "END:VCALENDAR")
	assert.Contains(t, output, "METHOD:PUBLISH")
	assert.Contains(t, output, "PRODID:-//Test//EN")
	assert.Contains(t, output, "X-WR-TIMEZONE:UTC")
	assert.Contains(t, output, "BEGIN:VEVENT")
	assert.Contains(t, output, "END:VEVENT")
	assert.Contains(t, output, "UID:event-1@demoparty.org")
	assert.Contains(t, output, "SUMMARY:Opening Ceremony")
	assert.Contains(t, output, "DESCRIPTION:The grand opening")
	assert.Contains(t, output, `LOCATION:Main Hall\, 123 Main St`)
	assert.Contains(t, output, "CATEGORIES:Ceremony")
	assert.Contains(t, output, "STATUS:CONFIRMED")
}

func TestIcalFormatter_Format_MultipleEntries(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)

	entries := domain.TimeTable{
		{
			ID:        1,
			Title:     "First",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        2,
			Title:     "Second",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(3 * time.Hour),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	f := NewIcalFormatter("", "", "")
	result, err := f.Format(entries)
	assert.NoError(t, err)

	output := string(result)
	assert.Equal(t, 2, strings.Count(output, "BEGIN:VEVENT"))
	assert.Contains(t, output, "UID:event-1@demoparty.org")
	assert.Contains(t, output, "UID:event-2@demoparty.org")
}

func TestIcalFormatter_Format_NoLocation(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)

	entries := domain.TimeTable{
		{
			ID:        1,
			Title:     "Virtual Event",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	f := NewIcalFormatter("", "", "")
	result, err := f.Format(entries)
	assert.NoError(t, err)

	output := string(result)
	assert.NotContains(t, output, "LOCATION:")
	assert.NotContains(t, output, "CATEGORIES:")
}

func TestIcalFormatter_Format_LocationWithDefaultAddress(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)

	entries := domain.TimeTable{
		{
			ID:        1,
			Title:     "Workshop",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			CreatedAt: now,
			UpdatedAt: now,
			Location: &domain.Location{
				Name: "Room B",
			},
		},
	}

	f := NewIcalFormatter("", "", "Venue Address")
	result, err := f.Format(entries)
	assert.NoError(t, err)

	output := string(result)
	assert.Contains(t, output, `LOCATION:Room B\, Venue Address`)
}

func TestIcalFormatter_Format_LocationWithEmptyAddress(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)
	emptyAddr := ""

	entries := domain.TimeTable{
		{
			ID:        1,
			Title:     "Workshop",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
			CreatedAt: now,
			UpdatedAt: now,
			Location: &domain.Location{
				Name:    "Room C",
				Address: &emptyAddr,
			},
		},
	}

	f := NewIcalFormatter("", "", "")
	result, err := f.Format(entries)
	assert.NoError(t, err)

	output := string(result)
	assert.Contains(t, output, "LOCATION:Room C")
	assert.NotContains(t, output, `LOCATION:Room C\,`)
}

func TestIcalFormatter_Format_Empty(t *testing.T) {
	f := NewIcalFormatter("", "", "")
	result, err := f.Format(domain.TimeTable{})
	assert.NoError(t, err)

	output := string(result)
	assert.Contains(t, output, "BEGIN:VCALENDAR")
	assert.Contains(t, output, "END:VCALENDAR")
	assert.NotContains(t, output, "BEGIN:VEVENT")
}
