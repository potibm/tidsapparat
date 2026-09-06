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
	data        []byte
	formatErr   error
	ext         string
	contentType string
}

func (m *mockFormatter) Format(entries domain.TimeTable) ([]byte, error) {
	return m.data, m.formatErr
}

func (m *mockFormatter) Extension() string {
	return m.ext
}

func (m *mockFormatter) ContentType() string {
	return m.contentType
}

type mockWriter struct {
	writeErr        error
	lastFile        string
	lastData        []byte
	lastContentType string
}

func (m *mockWriter) Write(_ context.Context, filename string, data []byte, contentType string) error {
	m.lastFile = filename
	m.lastData = data
	m.lastContentType = contentType

	return m.writeErr
}

func TestNewUniversalExporter(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	f := &mockFormatter{ext: ".ics", contentType: "text/calendar; charset=utf-8"}
	w := &mockWriter{}

	e := NewUniversalExporter("ical", "schedule", f, w, logger)

	assert.Equal(t, "ical", e.Name())
}

func TestUniversalExporter_Export_Success(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	f := &mockFormatter{data: []byte("ical-data"), ext: ".ics", contentType: "text/calendar; charset=utf-8"}
	w := &mockWriter{}

	e := NewUniversalExporter("ical", "schedule", f, w, logger)

	entries := domain.TimeTable{{ID: 1, Title: "Test"}}
	err := e.Export(context.Background(), entries)

	require.NoError(t, err)
	assert.Equal(t, "schedule.ics", w.lastFile)
	assert.Equal(t, []byte("ical-data"), w.lastData)
	assert.Equal(t, "text/calendar; charset=utf-8", w.lastContentType)
}

func TestUniversalExporter_Export_FormatError(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	f := &mockFormatter{
		formatErr:   errors.New("format failed"),
		ext:         ".ics",
		contentType: "text/calendar; charset=utf-8",
	}
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
	f := &mockFormatter{data: []byte("ical-data"), ext: ".ics", contentType: "text/calendar; charset=utf-8"}
	w := &mockWriter{writeErr: errors.New("write failed")}

	e := NewUniversalExporter("ical", "schedule", f, w, logger)

	entries := domain.TimeTable{{ID: 1, Title: "Test"}}
	err := e.Export(context.Background(), entries)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "write failed")
}
