package cmd

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/potibm/tidsapparat/internal/app/config"
	"github.com/potibm/tidsapparat/internal/app/exporter"
	"github.com/potibm/tidsapparat/internal/app/hub"
	"github.com/potibm/tidsapparat/internal/app/initializer"
	"github.com/potibm/tidsapparat/internal/app/repository"
	"github.com/potibm/tidsapparat/internal/app/services"
	store "github.com/potibm/tidsapparat/internal/app/store/gorm"
)

//go:embed assets
var staticFiles embed.FS

const (
	defaultDebounceDuration = 10 * time.Second
	otelEndpointFlagName    = "otel-endpoint"
	portFlagName            = "port"
)

var (
	port         int
	otelEndpoint string
)

func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Runs the HTTP server for the Tidsapparat application",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ensureAppInfrastructure()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := slog.Default()

			// =========================================================================
			// PHASE 1: Observability & Telemetry
			// =========================================================================
			shutdownFn, err := initializer.InitTelemetry(ctx, Cfg.App.OtelEndpoint, Cfg.App.Version)
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
				Port:              Cfg.App.Port,
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

	cmd.Flags().IntVarP(&port, portFlagName, "p", config.DefaultPort, "Set the port number for the server to listen on")
	_ = viper.BindPFlag("app.port", cmd.Flags().Lookup(portFlagName))

	cmd.Flags().
		StringVar(&otelEndpoint, otelEndpointFlagName, "", "Set the OpenTelemetry endpoint (e.g., localhost:4317)")
	_ = viper.BindPFlag("app.otel_endpoint", cmd.Flags().Lookup(otelEndpointFlagName))

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
