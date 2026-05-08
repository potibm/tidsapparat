package initializer

import (
	"context"
	"log/slog"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/potibm/billedapparat/internal/app/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

func TestBootstrapExporters_EmptyConfigs(t *testing.T) {
	exporters, err := BootstrapExporters(
		context.Background(),
		"1.0.0",
		config.PartyConfig{Timezone: "Europe/Berlin"},
		[]config.ExporterConfig{},
		nil,
		newTestLogger(),
	)

	require.NoError(t, err)
	assert.Empty(t, exporters)
}

func TestBootstrapExporters_UnknownType(t *testing.T) {
	configs := []config.ExporterConfig{
		{Name: "unknown-exporter", Type: "unknown", Destination: "file", Options: map[string]string{"dir": "/tmp"}},
	}

	exporters, err := BootstrapExporters(
		context.Background(),
		"1.0.0",
		config.PartyConfig{Timezone: "Europe/Berlin"},
		configs,
		nil,
		newTestLogger(),
	)

	require.NoError(t, err)
	assert.Empty(t, exporters)
}

func TestBootstrapExporters_UnknownDestination(t *testing.T) {
	configs := []config.ExporterConfig{
		{Name: "ical-ftp", Type: "ical", Destination: "ftp", Options: map[string]string{}},
	}

	exporters, err := BootstrapExporters(
		context.Background(),
		"1.0.0",
		config.PartyConfig{Timezone: "Europe/Berlin"},
		configs,
		nil,
		newTestLogger(),
	)

	require.NoError(t, err)
	assert.Empty(t, exporters)
}

func TestBootstrapExporters_IcalFileSuccess(t *testing.T) {
	configs := []config.ExporterConfig{
		{
			Name:        "ical-file",
			Type:        "ical",
			Destination: "file",
			Filename:    "schedule",
			Options:     map[string]string{"dir": "/tmp"},
		},
	}

	exporters, err := BootstrapExporters(
		context.Background(),
		"1.0.0",
		config.PartyConfig{Timezone: "Europe/Berlin", DefaultAddress: "Berlin"},
		configs,
		nil,
		newTestLogger(),
	)

	require.NoError(t, err)
	require.Len(t, exporters, 1)
	assert.Equal(t, "ical-file", exporters[0].Name())
}

func TestBootstrapExporters_FileMissingDir(t *testing.T) {
	configs := []config.ExporterConfig{
		{
			Name:        "ical-file",
			Type:        "ical",
			Destination: "file",
			Filename:    "schedule",
			Options:     map[string]string{},
		},
	}

	_, err := BootstrapExporters(
		context.Background(),
		"1.0.0",
		config.PartyConfig{Timezone: "Europe/Berlin"},
		configs,
		nil,
		newTestLogger(),
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires 'dir' option")
}

func TestBootstrapExporters_IcalS3Success(t *testing.T) {
	configs := []config.ExporterConfig{
		{
			Name:        "ical-s3",
			Type:        "ical",
			Destination: "s3",
			Filename:    "schedule",
			Options:     map[string]string{"bucket": "my-bucket"},
		},
	}

	mockS3 := &s3.Client{}

	exporters, err := BootstrapExporters(
		context.Background(),
		"1.0.0",
		config.PartyConfig{Timezone: "Europe/Berlin"},
		configs,
		mockS3,
		newTestLogger(),
	)

	require.NoError(t, err)
	require.Len(t, exporters, 1)
	assert.Equal(t, "ical-s3", exporters[0].Name())
}

func TestBootstrapExporters_S3MissingClient(t *testing.T) {
	configs := []config.ExporterConfig{
		{
			Name:        "ical-s3",
			Type:        "ical",
			Destination: "s3",
			Filename:    "schedule",
			Options:     map[string]string{"bucket": "my-bucket"},
		},
	}

	_, err := BootstrapExporters(
		context.Background(),
		"1.0.0",
		config.PartyConfig{Timezone: "Europe/Berlin"},
		configs,
		nil,
		newTestLogger(),
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "s3client is not configured")
}

func TestBootstrapExporters_S3MissingBucket(t *testing.T) {
	configs := []config.ExporterConfig{
		{
			Name:        "ical-s3",
			Type:        "ical",
			Destination: "s3",
			Filename:    "schedule",
			Options:     map[string]string{},
		},
	}

	mockS3 := &s3.Client{}

	_, err := BootstrapExporters(
		context.Background(),
		"1.0.0",
		config.PartyConfig{Timezone: "Europe/Berlin"},
		configs,
		mockS3,
		newTestLogger(),
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "requires 'bucket' option")
}

func TestBootstrapExporters_MixedConfigs(t *testing.T) {
	configs := []config.ExporterConfig{
		{
			Name:        "ical-file",
			Type:        "ical",
			Destination: "file",
			Filename:    "schedule",
			Options:     map[string]string{"dir": "/tmp"},
		},
		{
			Name:        "unknown-type",
			Type:        "csv",
			Destination: "file",
			Filename:    "data",
			Options:     map[string]string{"dir": "/tmp"},
		},
		{
			Name:        "unknown-dest",
			Type:        "ical",
			Destination: "ftp",
			Filename:    "calendar",
			Options:     map[string]string{},
		},
		{
			Name:        "ical-s3",
			Type:        "ical",
			Destination: "s3",
			Filename:    "schedule",
			Options:     map[string]string{"bucket": "my-bucket"},
		},
	}

	mockS3 := &s3.Client{}

	exporters, err := BootstrapExporters(
		context.Background(),
		"1.0.0",
		config.PartyConfig{Timezone: "Europe/Berlin"},
		configs,
		mockS3,
		newTestLogger(),
	)

	require.NoError(t, err)
	require.Len(t, exporters, 2)

	names := make([]string, len(exporters))
	for i, e := range exporters {
		names[i] = e.Name()
	}

	assert.Contains(t, names, "ical-file")
	assert.Contains(t, names, "ical-s3")
}

func TestBootstrapExporters_ProductIDIncludesVersion(t *testing.T) {
	configs := []config.ExporterConfig{
		{
			Name:        "ical-file",
			Type:        "ical",
			Destination: "file",
			Filename:    "schedule",
			Options:     map[string]string{"dir": "/tmp"},
		},
	}

	exporters, err := BootstrapExporters(
		context.Background(),
		"2.0.0",
		config.PartyConfig{Timezone: "Europe/Berlin"},
		configs,
		nil,
		newTestLogger(),
	)

	require.NoError(t, err)
	require.Len(t, exporters, 1)
	assert.Equal(t, "ical-file", exporters[0].Name())
}
