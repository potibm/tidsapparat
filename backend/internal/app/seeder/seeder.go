package seeder

import (
	"log/slog"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/potibm/billedapparat/internal/app/repository"
)

const (
	numberOfSponsorSlides     = 10
	numberOfSceneFriendSlides = 25
	numberOfNewsSlides        = 10
)

type Seeder struct {
	scheduleEntries []domain.ScheduleEntry
	currentID       int64
	repo            repository.ScheduleEntryRepository
}

func NewSeeder(repo repository.ScheduleEntryRepository) *Seeder {
	return &Seeder{
		scheduleEntries: []domain.ScheduleEntry{},
		currentID:       0,
		repo:            repo,
	}
}

func (s *Seeder) Run() error {
	slog.Info("Starting DB Purge & Seed...")

	_ = gofakeit.Seed(0)

	slog.Info("Seeding finished successfully", "total_slides", s.currentID)

	return nil
}
