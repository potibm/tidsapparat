package formatters

import (
	"fmt"

	ics "github.com/arran4/golang-ical"
	"github.com/potibm/billedapparat/internal/app/domain"
)

type IcalFormatter struct {
	ProductID     string
	Timezone      string
	DefaultAdress string
}

func NewIcalFormatter(productID, timezone, defaultAddress string) *IcalFormatter {
	if productID == "" {
		productID = "-//Tidsapparat//Timetable//EN"
	}

	if timezone == "" {
		timezone = "Europe/Berlin"
	}

	return &IcalFormatter{ProductID: productID, Timezone: timezone, DefaultAdress: defaultAddress}
}

func (f *IcalFormatter) Extension() string {
	return ".ics"
}

func (f *IcalFormatter) Format(entries domain.TimeTable) ([]byte, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetProductId(f.ProductID)

	cal.SetXWRTimezone(f.Timezone)

	for _, entry := range entries {
		uid := fmt.Sprintf("event-%d@demoparty.org", entry.ID)
		event := cal.AddEvent(uid)

		event.SetCreatedTime(entry.CreatedAt)
		event.SetDtStampTime(entry.UpdatedAt)
		event.SetStartAt(entry.StartTime)
		event.SetEndAt(entry.EndTime)

		event.SetSummary(entry.Title)
		event.SetDescription(entry.Description)

		if entry.Location != nil {
			locString := entry.Location.Name

			address := f.DefaultAdress
			if entry.Location.Address != nil && *entry.Location.Address != "" {
				address = *entry.Location.Address
			}

			if address != "" {
				locString = fmt.Sprintf("%s, %s", locString, address)
			}

			event.SetLocation(locString)
		}

		if entry.Category != nil {
			event.AddProperty(ics.ComponentPropertyCategories, entry.Category.Name)
		}

		event.AddProperty(ics.ComponentPropertySequence, fmt.Sprintf("%d", entry.UpdatedAt.Unix()))

		event.SetProperty(ics.ComponentPropertyStatus, string(ics.ObjectStatusConfirmed))
	}

	return []byte(cal.Serialize()), nil
}
