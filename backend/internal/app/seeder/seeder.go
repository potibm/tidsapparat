package seeder

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/potibm/tidsapparat/internal/app/domain"
	"github.com/potibm/tidsapparat/internal/app/repository"
)

const (
	locMainHall    = "Main Hall"
	locSeminarArea = "Seminar Area"
	locOutsideArea = "Outside Area"
	locClub        = "Club"
	catCompo       = "Compo"
	catEvent       = "Event"
	catSeminar     = "Seminar"
	catDeadline    = "Deadline"
	catGeneral     = "General"
	catLiveAct     = "Live Act"
	catDjSet       = "DJ Set"
)

type ScheduleEntrySeed struct {
	title, desc, url  string
	startDate         time.Time
	startTime         string
	endDate           *time.Time
	endTime, loc, cat string
}

type Seeder struct {
	categoryNamesToIDs map[string]int64
	locationNamesToIDs map[string]int64

	currentID         int64
	scheduleEntryRepo repository.ScheduleEntryRepository
	categoryRepo      repository.CategoryRepository
	locationRepo      repository.LocationRepository
}

func NewSeeder(
	scheduleEntryRepo repository.ScheduleEntryRepository,
	categoryRepo repository.CategoryRepository,
	locationRepo repository.LocationRepository,
) *Seeder {
	return &Seeder{
		categoryNamesToIDs: make(map[string]int64),
		locationNamesToIDs: make(map[string]int64),
		currentID:          0,
		scheduleEntryRepo:  scheduleEntryRepo,
		categoryRepo:       categoryRepo,
		locationRepo:       locationRepo,
	}
}

func (s *Seeder) Run() error {
	ctx := context.Background()

	slog.Info("Starting DB Seed...")

	_ = gofakeit.Seed(0)

	if err := s.seedLocations(ctx); err != nil {
		return fmt.Errorf("failed to seed locations: %w", err)
	}

	if err := s.seedCategories(ctx); err != nil {
		return fmt.Errorf("failed to seed categories: %w", err)
	}

	if err := s.seedSchedule(ctx); err != nil {
		return fmt.Errorf("failed to seed schedule: %w", err)
	}

	slog.Info("Seeding finished successfully")

	return nil
}

func (s *Seeder) seedLocations(ctx context.Context) error {
	address := createAddressString()

	locations := map[string]domain.Location{
		locMainHall:    createLocation(locMainHall, address),
		locSeminarArea: createLocation(locSeminarArea, address),
		locOutsideArea: createLocation(locOutsideArea, address),
		locClub:        createLocation(gofakeit.BuzzWord()+" Club", createAddressString()),
	}

	for locationName, location := range locations {
		slog.Info("Creating location", "name", location.Name, "address", *location.Address)

		if err := s.locationRepo.Create(ctx, &location); err != nil {
			return err
		} else {
			s.locationNamesToIDs[locationName] = location.ID
		}
	}

	return nil
}

func (s *Seeder) seedCategories(ctx context.Context) error {
	categories := map[string]domain.Category{
		catCompo:    createCategory(catCompo),
		catEvent:    createCategory(catEvent),
		catSeminar:  createCategory(catSeminar),
		catDeadline: createCategory(catDeadline),
		catGeneral:  createCategory(catGeneral),
		catLiveAct:  createCategory(catLiveAct),
		catDjSet:    createCategory(catDjSet),
	}

	for categoryName, category := range categories {
		slog.Info("Creating category", "name", category.Name, "color", category.Color)

		if err := s.categoryRepo.Create(ctx, &category); err != nil {
			return err
		} else {
			s.categoryNamesToIDs[categoryName] = category.ID
		}
	}

	return nil
}

func (s *Seeder) seedSchedule(ctx context.Context) error {
	friday := getWeekdayCurrentWeek(time.Friday)
	saturday := getWeekdayCurrentWeek(time.Saturday)
	sunday := getWeekdayCurrentWeek(time.Sunday)
	baseURL := "https://example.com/"
	scheduleEntries := make([]domain.ScheduleEntry, 0)

	entries := []ScheduleEntrySeed{
		{
			"Doors open",
			createDescription(),
			"",
			friday,
			"16:00",
			nil,
			"",
			locMainHall,
			catGeneral,
		},
		{
			"BBQ",
			createDescription(),
			"",
			friday,
			"17:00",
			nil,
			"18:00",
			locOutsideArea,
			catEvent,
		},
		{
			gofakeit.BuzzWord(),
			createDescription(),
			baseURL + "events#djset",
			friday,
			"17:00",
			nil,
			"18:00",
			locOutsideArea,
			catDjSet,
		},
		{"Opening Ceremony", "", "", friday, "19:00", nil, "19:30", locMainHall, catEvent},
		{"Demoshow", "", "", friday, "19:30", nil, "20:30", locMainHall, catEvent},
		{
			"Competition Show",
			createDescription(),
			baseURL + "events#competitionshow",
			friday,
			"20:30",
			nil,
			"22:00",
			locMainHall,
			catEvent,
		},
		{
			"Netlabel Night",
			createDescription(),
			baseURL + "events#netlabel",
			friday,
			"22:00",
			&saturday,
			"04:00",
			locClub,
			catEvent,
		},
		{
			createArtistName(),
			createDescription(),
			baseURL + "events#liveact",
			saturday,
			"00:30",
			nil,
			"02:00",
			locMainHall,
			catLiveAct,
		},
		{
			createArtistName(),
			createDescription(),
			baseURL + "events#liveact",
			saturday,
			"02:00",
			nil,
			"04:00",
			locMainHall,
			catLiveAct,
		},
		{
			"Deadline for Pixel graphics, One Screen ANSI/ASCII, Tracked Music, OGG/MP3 Music, " +
				"Animation, Freestyle Graphics, Alternative Platforms & Interactive",
			"",
			baseURL + "compos#deadlines",
			saturday,
			"00:00",
			nil,
			"00:00",
			"",
			catDeadline,
		},
		{
			"Deadline for Animation, Freestyle Graphics, 3D Scene & 4k Executable Graphics, " +
				"PC 4k-Intro, PC 64k-Intro, PC Demo",
			"",
			baseURL + "compos#deadlines",
			saturday,
			"12:00",
			nil,
			"12:00",
			"",
			catDeadline,
		},
		{
			createSeminarTitle(),
			createDescription(),
			baseURL + "seminars",
			saturday,
			"12:00",
			nil,
			"13:00",
			locSeminarArea,
			catSeminar,
		},
		{
			createSeminarTitle(),
			createDescription(),
			baseURL + "seminars",
			saturday,
			"13:00",
			nil,
			"14:00",
			locSeminarArea,
			catSeminar,
		},
		{
			"Tracked Music & One Screen ANSI/ASCII",
			createDescription(),
			baseURL + "compos",
			saturday,
			"14:00",
			nil,
			"15:00",
			locMainHall,
			catCompo,
		},
		{
			"BBQ",
			createDescription(),
			"",
			saturday,
			"15:00",
			nil,
			"16:00",
			locOutsideArea,
			catEvent,
		},
		{
			createDJName(),
			createDescription(),
			"",
			saturday,
			"15:00",
			nil,
			"16:00",
			locOutsideArea,
			catDjSet,
		},
		{
			createSeminarTitle(),
			createDescription(),
			baseURL + "seminars",
			saturday,
			"16:00",
			nil,
			"17:00",
			locSeminarArea,
			catSeminar,
		},
		{
			"OGG/MP3 Music",
			createDescription(),
			baseURL + "compos",
			saturday,
			"17:30",
			nil,
			"18:30",
			locMainHall,
			catCompo,
		},
		{
			"Pixel Graphics, Freestyle Graphics, Animation, 3D Scene & 4k Executable Graphics",
			createDescription(),
			baseURL + "compos",
			saturday,
			"19:00",
			nil,
			"20:30",
			locMainHall,
			catCompo,
		},
		{
			"Interactive & Alternative Platform",
			createDescription(),
			baseURL + "compos",
			saturday,
			"21:30",
			nil,
			"22:30",
			locMainHall,
			catCompo,
		},
		{
			"PC 4k-Intro, PC 64k-Intro & PC Demo",
			createDescription(),
			baseURL + "compos",
			saturday,
			"23:00",
			&sunday,
			"00:30",
			locMainHall,
			catCompo,
		},
		{
			createDJName(),
			createDescription(),
			"",
			sunday,
			"01:00",
			nil,
			"02:30",
			locOutsideArea,
			catDjSet,
		},
		{
			createDJName(),
			createDescription(),
			"",
			sunday,
			"02:30",
			nil,
			"04:30",
			locOutsideArea,
			catDjSet,
		},
		{"Voting Deadline", createDescription(), "", sunday, "08:00", nil, "08:00", "", catDeadline},
		{
			createDJName(),
			"",
			"",
			sunday,
			"11:00",
			nil,
			"12:45",
			locOutsideArea,
			catDjSet,
		},
		{
			"Prizegiving",
			createDescription(),
			"",
			sunday,
			"13:00",
			nil,
			"14:00",
			locMainHall,
			catEvent,
		},
	}

	for _, e := range entries {
		start := getDayAtTime(e.startDate, e.startTime)

		endDate := e.startDate
		if e.endDate != nil {
			endDate = *e.endDate
		}

		endTime := e.endTime
		if endTime == "" {
			endTime = e.startTime
		}

		end := getDayAtTime(endDate, endTime)
		scheduleEntry := s.createScheduleEntry(e.title, e.desc, e.url, start, end, e.loc, e.cat)
		scheduleEntries = append(scheduleEntries, scheduleEntry)
	}

	for _, entry := range scheduleEntries {
		slog.Info("Creating schedule entry", "name", entry.Title)

		if err := s.scheduleEntryRepo.Save(ctx, &entry); err != nil {
			return err
		}
	}

	return nil
}

func createLocation(name, address string) domain.Location {
	if address == "" {
		address = gofakeit.Address().Address
	}

	return domain.Location{
		Name:    name,
		Address: &address,
	}
}

func createCategory(name string) domain.Category {
	return domain.Category{
		Name:  name,
		Color: gofakeit.HexColor(),
	}
}

func (s *Seeder) findLocationIDByName(name string) *int64 {
	if id, exists := s.locationNamesToIDs[name]; exists {
		res := id

		return &res
	}

	return nil
}

func (s *Seeder) findCategoryIDByName(name string) *int64 {
	if id, exists := s.categoryNamesToIDs[name]; exists {
		res := id

		return &res
	}

	return nil
}

func (s *Seeder) createScheduleEntry(
	title, description, externalURL string,
	startTime, endTime time.Time,
	locationName string,
	categoryName string,
) domain.ScheduleEntry {
	var (
		categoryPointer *int64 = nil
		locationPointer *int64 = nil
	)

	if categoryName != "" {
		categoryPointer = s.findCategoryIDByName(categoryName)
		if categoryPointer == nil {
			slog.Warn("Category not found", "name", categoryName)
		}
	}

	if locationName != "" {
		locationPointer = s.findLocationIDByName(locationName)
		if locationPointer == nil {
			slog.Warn("Location not found", "name", locationName)
		}
	}

	scheduleEntry := domain.ScheduleEntry{
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		ExternalURL: externalURL,
		LocationID:  locationPointer,
		CategoryID:  categoryPointer,
	}

	return scheduleEntry
}

func createAddressString() string {
	return gofakeit.Address().Address
}

func getWeekdayRelative(now time.Time, weekday time.Weekday) time.Time {
	currentWeekday := int(now.Weekday())
	targetWeekday := int(weekday)

	if targetWeekday == int(time.Sunday) {
		targetWeekday = 7
	}

	if currentWeekday == int(time.Sunday) {
		currentWeekday = 7
	}

	daysUntilWeekday := targetWeekday - currentWeekday

	return now.AddDate(0, 0, daysUntilWeekday)
}

func getWeekdayCurrentWeek(weekday time.Weekday) time.Time {
	return getWeekdayRelative(time.Now(), weekday)
}

func getDayAtTime(date time.Time, timeStr string) time.Time {
	var hour, minute int

	_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	if err != nil {
		slog.Error("Failed to parse time string", "time_str", timeStr, "error", err)

		return date
	}

	return time.Date(
		date.Year(), date.Month(), date.Day(),
		hour, minute, 0, 0,
		date.Location(),
	)
}

func createDJName() string {
	return fmt.Sprintf("DJ %s", gofakeit.Gamertag())
}

func createArtistName() string {
	return gofakeit.AppName()
}

func createSeminarTitle() string {
	return gofakeit.Verb() + " the " + gofakeit.Noun()
}

func createDescription() string {
	return gofakeit.Sentence()
}
