package exporter

import (
	"context"

	"github.com/potibm/billedapparat/internal/app/domain"
)

type Formatter interface {
	Format(entries domain.TimeTable) ([]byte, error)
	Extension() string // e.g. ".ics" or ".json"
}

type Writer interface {
	Write(ctx context.Context, filename string, data []byte) error
}
