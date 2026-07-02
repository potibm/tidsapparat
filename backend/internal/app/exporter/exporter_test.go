package exporter

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/potibm/tidsapparat/internal/app/domain"
	"github.com/stretchr/testify/assert"
)

// --- Mock implementations ---

type mockAllPreloader struct {
	result domain.TimeTable
	err    error
}

func (m *mockAllPreloader) GetAllPreloaded(ctx context.Context) (domain.TimeTable, error) {
	return m.result, m.err
}

type mockExporter struct {
	name        string
	exportErr   error
	exportCount int
	lastEntries domain.TimeTable
}

func (m *mockExporter) Name() string { return m.name }

func (m *mockExporter) Export(ctx context.Context, entries domain.TimeTable) error {
	m.exportCount++
	m.lastEntries = entries

	return m.exportErr
}

func TestNewManager(t *testing.T) {
	preloader := &mockAllPreloader{}
	m := NewManager(preloader, 100*time.Millisecond)

	assert.NotNil(t, m)
	assert.Empty(t, m.exporters)
	assert.Equal(t, 100*time.Millisecond, m.debounceTime)
}

func TestManager_Register(t *testing.T) {
	preloader := &mockAllPreloader{}
	m := NewManager(preloader, 100*time.Millisecond)

	e1 := &mockExporter{name: "e1"}
	e2 := &mockExporter{name: "e2"}

	m.Register(e1)
	m.Register(e2)

	assert.Len(t, m.exporters, 2)
}

func TestManager_Ping_Debounce(t *testing.T) {
	preloader := &mockAllPreloader{
		result: domain.TimeTable{
			{ID: 1, Title: "Visible"},
		},
	}
	m := NewManager(preloader, 50*time.Millisecond)

	e1 := &mockExporter{name: "e1"}
	m.Register(e1)

	// Multiple pings should be debounced
	m.Ping()
	m.Ping()
	m.Ping()

	// Wait for debounce to fire
	time.Sleep(150 * time.Millisecond)

	// Should only run once due to debounce
	assert.Equal(t, 1, e1.exportCount)
}

func TestManager_RunAll_FiltersHidden(t *testing.T) {
	preloader := &mockAllPreloader{
		result: domain.TimeTable{
			{ID: 1, Title: "Visible", Hidden: false},
			{ID: 2, Title: "Hidden", Hidden: true},
			{ID: 3, Title: "Also Visible", Hidden: false},
		},
	}
	m := NewManager(preloader, 0)

	e1 := &mockExporter{name: "e1"}
	m.Register(e1)

	m.RunAll()

	assert.Equal(t, 2, len(e1.lastEntries))
	assert.Equal(t, int64(1), e1.lastEntries[0].ID)
	assert.Equal(t, int64(3), e1.lastEntries[1].ID)
}

func TestManager_RunAll_PreloaderError(t *testing.T) {
	preloader := &mockAllPreloader{err: errors.New("db down")}
	m := NewManager(preloader, 0)

	e1 := &mockExporter{name: "e1"}
	m.Register(e1)

	m.RunAll()

	assert.Equal(t, 0, e1.exportCount)
}

func TestManager_RunAll_ExporterError(t *testing.T) {
	preloader := &mockAllPreloader{
		result: domain.TimeTable{{ID: 1, Title: "Event"}},
	}
	m := NewManager(preloader, 0)

	e1 := &mockExporter{name: "failer", exportErr: errors.New("export failed")}
	e2 := &mockExporter{name: "success"}

	m.Register(e1)
	m.Register(e2)

	m.RunAll()

	// Both should be called, even if one fails
	assert.Equal(t, 1, e1.exportCount)
	assert.Equal(t, 1, e2.exportCount)
}

func TestManager_RunAll_ConcurrentExporters(t *testing.T) {
	preloader := &mockAllPreloader{
		result: domain.TimeTable{{ID: 1, Title: "Event"}},
	}
	m := NewManager(preloader, 0)

	exporters := make([]*mockExporter, 5)
	for i := range exporters {
		exporters[i] = &mockExporter{name: "e" + string(rune('0'+i))}
		m.Register(exporters[i])
	}

	m.RunAll()

	for _, e := range exporters {
		assert.Equal(t, 1, e.exportCount, "exporter %s should be called once", e.name)
	}
}
