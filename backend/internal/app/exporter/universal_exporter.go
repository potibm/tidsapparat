package exporter

import (
	"context"
	"log/slog"

	"github.com/potibm/billedapparat/internal/app/domain"
)

type UniversalExporter struct {
	name      string
	formatter Formatter
	writer    Writer
	filename  string
	logger    *slog.Logger
}

func NewUniversalExporter(name, filename string, f Formatter, w Writer, l *slog.Logger) *UniversalExporter {
	return &UniversalExporter{
		name:      name,
		formatter: f,
		writer:    w,
		filename:  filename,
		logger:    l.With("exporter", name),
	}
}

func (e *UniversalExporter) Name() string { return e.name }

func (e *UniversalExporter) Export(ctx context.Context, entries domain.TimeTable) error {
	e.logger.Info("Starting export run", "entry_count", len(entries))

	data, err := e.formatter.Format(entries)
	if err != nil {
		e.logger.Error("Formatting failed", "error", err)

		return err
	}

	fullFilename := e.filename + e.formatter.Extension()

	err = e.writer.Write(ctx, fullFilename, data)
	if err != nil {
		e.logger.Error("Writing failed", "filename", fullFilename, "error", err)

		return err
	}

	e.logger.Info("Export successful", "filename", fullFilename, "bytes", len(data))

	return nil
}
