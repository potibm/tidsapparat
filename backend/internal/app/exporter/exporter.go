package exporter

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/potibm/billedapparat/internal/app/domain"
)

const defaultTimeout = 30 * time.Second

type Exporter interface {
	Name() string
	Export(ctx context.Context, timetable domain.TimeTable) error
}

type Manager struct {
	exporters    []Exporter
	db           AllPreloader // Interface to fetch fresh data from the database
	debounceTime time.Duration
	timer        *time.Timer
	mu           sync.Mutex
	logger       *slog.Logger
}

type AllPreloader interface {
	GetAllPreloaded(ctx context.Context) (domain.TimeTable, error)
}

func NewManager(source AllPreloader, debounce time.Duration) *Manager {
	logger := slog.Default()

	return &Manager{
		exporters:    []Exporter{},
		db:           source,
		debounceTime: debounce,
		logger:       logger.With("component", "Exporter"),
	}
}

func (m *Manager) Register(e Exporter) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.exporters = append(m.exporters, e)
}

func (m *Manager) Ping() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.timer != nil {
		m.timer.Stop()
	}

	m.timer = time.AfterFunc(m.debounceTime, func() {
		m.RunAll()
	})
}

func (m *Manager) RunAll() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	timetable, err := m.db.GetAllPreloaded(ctx)
	if err != nil {
		m.logger.Error("Error fetching timetable", "error", err)

		return
	}

	var wg sync.WaitGroup
	for _, e := range m.exporters {
		wg.Add(1)

		go func(exp Exporter) {
			defer wg.Done()

			m.logger.Info("Starting", "exporter", exp.Name())

			if err := exp.Export(ctx, timetable); err != nil {
				m.logger.Error("Failed", "exporter", exp.Name(), "error", err)
			} else {
				m.logger.Info("Finished successfully", "exporter", exp.Name())
			}
		}(e)
	}

	wg.Wait()
}
