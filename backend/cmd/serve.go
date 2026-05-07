package cmd

import (
	"embed"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"

	"github.com/potibm/billedapparat/internal/app/exporter"
	"github.com/potibm/billedapparat/internal/app/hub"
	"github.com/potibm/billedapparat/internal/app/initializer"
	"github.com/potibm/billedapparat/internal/app/services"
	store "github.com/potibm/billedapparat/internal/app/store/gorm"
)

//go:embed assets
var staticFiles embed.FS

const (
	defaultPort = 3200
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
			// 1. Context
			ctx := cmd.Context()

			// 2. Initialize Telemetry
			shutdownFn, err := initializer.InitTelemetry(ctx, otelEndpoint, Cfg.App.Version)
			if err != nil {
				return fmt.Errorf("failed to initialize telemetry: %w", err)
			}

			if shutdownFn != nil {
				defer shutdownFn()
			}

			dbStore, err := store.NewSqliteStore(Cfg.App.DbFilename)
			if err != nil {
				return fmt.Errorf("database error: %w", err)
			}

			defer func() {
				if err := dbStore.Close(); err != nil {
					slog.Error("failed to close database", "error", err)
				}
			}()

			redisClient := initializer.InitializeRedis(Cfg.App.RedisURL)

			scheduleRepo := dbStore.NewScheduleEntryRepository()

			slog.Info("Performing initial boot export...")

			exportMgr := exporter.NewManager(scheduleRepo, 10*time.Second)
			exportMgr.Register(exporter.NewLogExporter())

			go exportMgr.RunAll()

			slog.Info("Performing initial boot sync...")

			eventHub := services.NewEventHub(exportMgr, redisClient, scheduleRepo)
			eventHub.PublishFullSync(ctx)

			// 4. Initialize external services (Sentry)
			initializer.InitializeSentry(Cfg.Sentry)

			// 5. Start the server
			server, err := hub.NewServer(hub.Config{
				Port:              port,
				StaticFiles:       staticFiles,
				ScheduleEntryRepo: scheduleRepo,
				CategoryRepo:      dbStore.NewCategoryRepository(),
				LocationRepo:      dbStore.NewLocationRepository(),
				EventHub:          eventHub,
				ExporterManager:   exportMgr,
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
