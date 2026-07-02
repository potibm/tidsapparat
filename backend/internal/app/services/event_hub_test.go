package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/potibm/protokolapparat/pkg/schedule"
	"github.com/potibm/tidsapparat/internal/app/domain"
	"github.com/potibm/tidsapparat/internal/app/exporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock ScheduleSource ---

type mockScheduleSource struct {
	getByIDResult         *domain.ScheduleEntry
	getByIDErr            error
	getAllPreloadedResult domain.TimeTable
	getAllPreloadedErr    error
	getByCategoryIDResult []domain.ScheduleEntry
	getByLocationIDResult []domain.ScheduleEntry
}

func (m *mockScheduleSource) GetByID(ctx context.Context, id int64) (*domain.ScheduleEntry, error) {
	return m.getByIDResult, m.getByIDErr
}

func (m *mockScheduleSource) GetAllPreloaded(ctx context.Context) (domain.TimeTable, error) {
	return m.getAllPreloadedResult, m.getAllPreloadedErr
}

func (m *mockScheduleSource) GetByCategoryID(ctx context.Context, categoryID int64) ([]domain.ScheduleEntry, error) {
	return m.getByCategoryIDResult, nil
}

func (m *mockScheduleSource) GetByLocationID(ctx context.Context, locationID int64) ([]domain.ScheduleEntry, error) {
	return m.getByLocationIDResult, nil
}

// --- Mock Exporter ---

type mockTestExporter struct {
	exportCount int
	lastEntries domain.TimeTable
}

func (m *mockTestExporter) Name() string { return "mock" }

func (m *mockTestExporter) Export(ctx context.Context, entries domain.TimeTable) error {
	m.exportCount++
	m.lastEntries = entries

	return nil
}

func newTestEventHub(repo ScheduleSource, debounce time.Duration) *EventHub {
	expMgr := exporter.NewManager(repo, debounce)

	return NewEventHub(expMgr, nil, repo)
}

func TestNewEventHub(t *testing.T) {
	repo := &mockScheduleSource{}
	hub := newTestEventHub(repo, 0)

	assert.NotNil(t, hub)
	assert.NotNil(t, hub.logger)
}

func TestEventHub_PublishCreate_Success(t *testing.T) {
	repo := &mockScheduleSource{
		getByIDResult: &domain.ScheduleEntry{
			ID:        1,
			Title:     "Test",
			StartTime: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
		},
	}
	hub := newTestEventHub(repo, 0)

	// With nil redis, send returns early without panic
	hub.PublishCreate(context.Background(), 1)
	assert.NotNil(t, hub)
}

func TestEventHub_PublishCreate_RepoError(t *testing.T) {
	repo := &mockScheduleSource{
		getByIDErr: errors.New("not found"),
	}
	hub := newTestEventHub(repo, 0)

	// Should not panic when repo returns error
	hub.PublishCreate(context.Background(), 99)
	assert.NotNil(t, hub)
}

func TestEventHub_PublishUpdate_Success(t *testing.T) {
	repo := &mockScheduleSource{
		getByIDResult: &domain.ScheduleEntry{
			ID:        2,
			Title:     "Updated",
			StartTime: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
		},
	}
	hub := newTestEventHub(repo, 0)

	hub.PublishUpdate(context.Background(), 2)
	assert.NotNil(t, hub)
}

func TestEventHub_PublishUpdate_RepoError(t *testing.T) {
	repo := &mockScheduleSource{
		getByIDErr: errors.New("not found"),
	}
	hub := newTestEventHub(repo, 0)

	hub.PublishUpdate(context.Background(), 99)
	assert.NotNil(t, hub)
}

func TestEventHub_PublishDelete(t *testing.T) {
	repo := &mockScheduleSource{}
	hub := newTestEventHub(repo, 0)

	// With nil redis, send returns early without panic
	hub.PublishDelete(context.Background(), 3)
	assert.NotNil(t, hub)
}

func TestEventHub_PublishFullSync_NilRedis(t *testing.T) {
	repo := &mockScheduleSource{}
	hub := newTestEventHub(repo, 0)

	// Should return early without error when redis is nil
	hub.PublishFullSync(context.Background())
	assert.NotNil(t, hub)
}

func TestEventHub_PublishFullSync_RepoError(t *testing.T) {
	repo := &mockScheduleSource{
		getAllPreloadedErr: errors.New("db error"),
	}
	hub := newTestEventHub(repo, 0)

	hub.PublishFullSync(context.Background())
	assert.NotNil(t, hub)
}

func TestEventHub_SyncCategoryUpdate(t *testing.T) {
	now := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	repo := &mockScheduleSource{
		getByCategoryIDResult: []domain.ScheduleEntry{
			{ID: 1, Title: "Entry 1", StartTime: now, EndTime: now.Add(time.Hour)},
			{ID: 2, Title: "Entry 2", StartTime: now, EndTime: now.Add(time.Hour)},
		},
		getByIDResult: &domain.ScheduleEntry{
			ID:        1,
			Title:     "Entry 1",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
		},
	}
	exp := &mockTestExporter{}
	expMgr := exporter.NewManager(repo, 0)
	expMgr.Register(exp)
	hub := NewEventHub(expMgr, nil, repo)

	hub.SyncCategoryUpdate(context.Background(), 5)

	// SyncCategoryUpdate calls exporter.Ping which triggers async RunAll
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 1, exp.exportCount)
}

func TestEventHub_SyncLocationUpdate(t *testing.T) {
	now := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	repo := &mockScheduleSource{
		getByLocationIDResult: []domain.ScheduleEntry{
			{ID: 3, Title: "Entry 3", StartTime: now, EndTime: now.Add(time.Hour)},
		},
		getByIDResult: &domain.ScheduleEntry{
			ID:        3,
			Title:     "Entry 3",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
		},
	}
	exp := &mockTestExporter{}
	expMgr := exporter.NewManager(repo, 0)
	expMgr.Register(exp)
	hub := NewEventHub(expMgr, nil, repo)

	hub.SyncLocationUpdate(context.Background(), 7)

	// SyncLocationUpdate calls exporter.Ping which triggers async RunAll
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, 1, exp.exportCount)
}

func TestEventHub_getProtocolEntry(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)
	repo := &mockScheduleSource{
		getByIDResult: &domain.ScheduleEntry{
			ID:        1,
			Title:     "Test",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
		},
	}
	hub := newTestEventHub(repo, 0)

	entry, err := hub.getProtocolEntry(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, schedule.Entry{
		ID:        1,
		Title:     "Test",
		StartTime: "2024-06-15T10:00:00Z",
		EndTime:   "2024-06-15T11:00:00Z",
		IsHidden:  false,
	}, entry)
}

func TestEventHub_getProtocolEntry_Error(t *testing.T) {
	repo := &mockScheduleSource{
		getByIDErr: errors.New("not found"),
	}
	hub := newTestEventHub(repo, 0)

	_, err := hub.getProtocolEntry(context.Background(), 99)
	assert.Error(t, err)
}
