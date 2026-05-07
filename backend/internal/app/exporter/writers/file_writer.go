package writers

import (
	"context"
	"os"
	"path/filepath"
)

const (
	defaultFileMode = 0o600
	defaultDirMode  = 0o755
)

type FileWriter struct {
	BaseDir string
}

func (w *FileWriter) Write(ctx context.Context, filename string, data []byte) error {
	path := filepath.Join(w.BaseDir, filename)
	if err := os.MkdirAll(w.BaseDir, defaultDirMode); err != nil {
		return err
	}

	return os.WriteFile(path, data, defaultFileMode)
}
