package cmd

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"

	"github.com/potibm/billedapparat/internal/app/exporter"
	"github.com/potibm/billedapparat/internal/app/hub"
	"github.com/potibm/billedapparat/internal/app/initializer"
	"github.com/potibm/billedapparat/internal/app/repository"
	"github.com/potibm/billedapparat/internal/app/services"
	store "github.com/potibm/billedapparat/internal/app/store/gorm"
)

//go:embed assets
var staticFiles embed.FS

const (
	defaultPort             = 3200
	defaultDebounceDuration = 10 * time.Second
)

var (
	port         int
	otelEndpoint string
)

func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Runs the HTTP server for the Billedapparat application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ensureAppInfrastructure()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := slog.Default()

			// =========================================================================
			// PHASE 1: Observability & Telemetry
			// =========================================================================
			shutdownFn, err := initializer.InitTelemetry(ctx, otelEndpoint, Cfg.App.Version)
			if err != nil {
				return fmt.Errorf("failed to initialize telemetry: %w", err)
			}

			if shutdownFn != nil {
				defer shutdownFn()
			}

			initializer.InitializeSentry(Cfg.Sentry)

			// =========================================================================
			// PHASE 2: Core Infrastructure (DB & Redis)
			// =========================================================================
			dbStore, err := store.NewSqliteStore(Cfg.App.DbFilename)
			if err != nil {
				return fmt.Errorf("database error: %w", err)
			}

			defer func() {
				if err := dbStore.Close(); err != nil {
					logger.Error("failed to close database", "error", err)
				}
			}()

			scheduleRepo := dbStore.NewScheduleEntryRepository()

			redisClient := initializer.InitializeRedis(Cfg.App.RedisURL)

			// =========================================================================
			// PHASE 3: Background Tasks (Exporters & Hub)
			// =========================================================================

			exportMgr, err := setupExportManager(ctx, scheduleRepo, logger)
			if err != nil {
				return err
			}

			// Trigger initial export and sync
			slog.Info("Performing initial boot export...")

			go exportMgr.RunAll()

			slog.Info("Performing initial boot sync...")

			eventHub := services.NewEventHub(exportMgr, redisClient, scheduleRepo)
			eventHub.PublishFullSync(ctx)

			// =========================================================================
			// PHASE 4: API Server
			// =========================================================================

			server, err := hub.NewServer(hub.Config{
				Port:              port,
				StaticFiles:       staticFiles,
				ScheduleEntryRepo: scheduleRepo,
				CategoryRepo:      dbStore.NewCategoryRepository(),
				LocationRepo:      dbStore.NewLocationRepository(),
				EventHub:          eventHub,
				Cfg:               Cfg,
			})
			if err != nil {
				return fmt.Errorf("failed to initialize server: %w", err)
			}

			return server.Run(ctx)
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", defaultPort, "Set the port number for the server to listen on")
	cmd.Flags().
		StringVar(&otelEndpoint, "otel-endpoint", "", "Set the OpenTelemetry endpoint (e.g., localhost:4317)")

	return cmd
}

func setupExportManager(
	ctx context.Context,
	repo repository.ScheduleEntryRepository,
	logger *slog.Logger,
) (*exporter.Manager, error) {
	exportMgr := exporter.NewManager(repo, defaultDebounceDuration)

	// Always register default log exporter for debugging purposes
	exportMgr.Register(exporter.NewLogExporter())

	// Initialize s3 client if configured
	s3Client, err := initializer.InitializeS3Client(ctx, Cfg.S3Client, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 client: %w", err)
	}

	// Create exporters based on config
	exporters, err := initializer.BootstrapExporters(
		ctx,
		Cfg.App.Version,
		Cfg.Party,
		Cfg.Exporter,
		s3Client,
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bootstrap exporters: %w", err)
	}

	// Register exporters with the manager
	for _, exp := range exporters {
		exportMgr.Register(exp)
	}

	return exportMgr, nil
}
