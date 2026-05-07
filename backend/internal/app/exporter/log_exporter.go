package exporter

import (
	"context"
	"log/slog"

	"github.com/potibm/billedapparat/internal/app/domain"
)

type LogExporter struct{}

func NewLogExporter() *LogExporter {
	return &LogExporter{}
}

func (e *LogExporter) Name() string {
	return "LogExporter"
}

func (e *LogExporter) Export(ctx context.Context, entries domain.TimeTable) error {
	slog.Info("Exporting timetable to log", "count", len(entries))

	for _, entry := range entries {
		locName := "N/A"
		if entry.Location != nil {
			locName = entry.Location.Name
		}

		slog.Debug("Schedule Entry",
			"title", entry.Title,
			"location", locName,
			"start", entry.StartTime)
	}

	return nil
}
