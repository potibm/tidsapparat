package exporter

import (
	"context"

	"github.com/potibm/tidsapparat/internal/app/domain"
)

type Formatter interface {
	Format(entries domain.TimeTable) ([]byte, error)
	Extension() string   // e.g. ".ics" or ".json"
	ContentType() string // e.g. "text/calendar; charset=utf-8"
}

type Writer interface {
	Write(ctx context.Context, filename string, data []byte, contentType string) error
}
