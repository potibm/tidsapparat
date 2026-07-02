package exporter

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/potibm/tidsapparat/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---

type mockFormatter struct {
	data      []byte
	formatErr error
	ext       string
}

func (m *mockFormatter) Format(entries domain.TimeTable) ([]byte, error) {
	return m.data, m.formatErr
}

func (m *mockFormatter) Extension() string {
	return m.ext
}

type mockWriter struct {
	writeErr error
	lastFile string
	lastData []byte
}

func (m *mockWriter) Write(ctx context.Context, filename string, data []byte) error {
	m.lastFile = filename
	m.lastData = data

	return m.writeErr
}

func TestNewUniversalExporter(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	f := &mockFormatter{ext: ".ics"}
	w := &mockWriter{}

	e := NewUniversalExporter("ical", "schedule", f, w, logger)

	assert.Equal(t, "ical", e.Name())
}

func TestUniversalExporter_Export_Success(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	f := &mockFormatter{data: []byte("ical-data"), ext: ".ics"}
	w := &mockWriter{}

	e := NewUniversalExporter("ical", "schedule", f, w, logger)

	entries := domain.TimeTable{{ID: 1, Title: "Test"}}
	err := e.Export(context.Background(), entries)

	require.NoError(t, err)
	assert.Equal(t, "schedule.ics", w.lastFile)
	assert.Equal(t, []byte("ical-data"), w.lastData)
}

func TestUniversalExporter_Export_FormatError(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	f := &mockFormatter{formatErr: errors.New("format failed"), ext: ".ics"}
	w := &mockWriter{}

	e := NewUniversalExporter("ical", "schedule", f, w, logger)

	entries := domain.TimeTable{{ID: 1, Title: "Test"}}
	err := e.Export(context.Background(), entries)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "format failed")
	assert.Empty(t, w.lastFile)
}

func TestUniversalExporter_Export_WriteError(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	f := &mockFormatter{data: []byte("ical-data"), ext: ".ics"}
	w := &mockWriter{writeErr: errors.New("write failed")}

	e := NewUniversalExporter("ical", "schedule", f, w, logger)

	entries := domain.TimeTable{{ID: 1, Title: "Test"}}
	err := e.Export(context.Background(), entries)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "write failed")
}
