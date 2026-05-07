package hub

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/potibm/billedapparat/internal/app/config"
	"github.com/potibm/billedapparat/internal/app/repository"
	"github.com/potibm/billedapparat/internal/app/services"
	sloggin "github.com/samber/slog-gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	defaultShutdownTimeout    = 5 * time.Second
	defaultReadHeaderTimeout  = 3 * time.Second
	pathScheduleEntries       = "/schedule-entries"
	pathScheduleEntriesWithID = "/schedule-entries/:id"
	pathCategories            = "/categories"
	pathCategoriesWithID      = "/categories/:id"
	pathLocations             = "/locations"
	pathLocationsWithID       = "/locations/:id"
)

type Config struct {
	Port              int
	StaticFiles       embed.FS
	ScheduleEntryRepo repository.ScheduleEntryRepository
	CategoryRepo      repository.CategoryRepository
	LocationRepo      repository.LocationRepository
	EventHub          *services.EventHub
	Cfg               config.Config
}

type Server struct {
	port              int
	staticFiles       embed.FS
	eventHub          *services.EventHub
	scheduleEntryRepo repository.ScheduleEntryRepository
	categoryRepo      repository.CategoryRepository
	locationRepo      repository.LocationRepository
	cfg               config.Config
	logger            *slog.Logger
}

func NewServer(cfg Config) (*Server, error) {
	logger := slog.Default()

	if cfg.EventHub == nil {
		return nil, fmt.Errorf("event hub is nil")
	}

	return &Server{
		port:              cfg.Port,
		staticFiles:       cfg.StaticFiles,
		scheduleEntryRepo: cfg.ScheduleEntryRepo,
		categoryRepo:      cfg.CategoryRepo,
		locationRepo:      cfg.LocationRepo,
		cfg:               cfg.Cfg,
		eventHub:          cfg.EventHub,
		logger:            logger.With("component", "HubServer"),
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	router, err := s.setupRouter()
	if err != nil {
		return fmt.Errorf("setup router: %w", err)
	}

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(s.port),
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		Handler:           router,
	}

	// Start server in Goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen", "error", err)
		}
	}()

	slog.Info("Server is up", "port", s.port)

	<-ctx.Done()
	slog.Info("Shutting down server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}

func (s *Server) setupRouter() (*gin.Engine, error) {
	gin.SetMode(s.cfg.App.GinMode)

	r := gin.New()

	r.Use(
		// middleware.ErrorHandlingMiddleware(),
		gin.Recovery(),
		sentrygin.New(sentrygin.Options{Repanic: false}),
		sloggin.New(slog.Default()),
		otelgin.Middleware(config.OtelBackendServiceName),
	)
	s.registerCorsMiddleware(r)

	r.Static("/media", "./data/media")
	r.Static("/style", "./data/style")

	folder, err := static.EmbedFolder(s.staticFiles, "assets")
	if err != nil {
		return nil, fmt.Errorf("create embedded folder: %w", err)
	}

	r.Use(static.Serve("/", folder))

	api := r.Group("/api")
	api.GET("/config", s.handleGetPublicConfig)

	admin := r.Group("/api/admin")
	admin.GET(pathScheduleEntries, s.listScheduleEntries)
	admin.POST(pathScheduleEntries, s.createScheduleEntry)
	admin.GET(pathScheduleEntriesWithID, s.getScheduleEntry)
	admin.PUT(pathScheduleEntriesWithID, s.updateScheduleEntry)
	admin.DELETE(pathScheduleEntriesWithID, s.deleteScheduleEntry)

	admin.GET(pathCategories, s.listCategories)
	admin.POST(pathCategories, s.createCategory)
	admin.GET(pathCategoriesWithID, s.getCategory)
	admin.PUT(pathCategoriesWithID, s.updateCategory)
	admin.DELETE(pathCategoriesWithID, s.deleteCategory)

	admin.GET(pathLocations, s.listLocations)
	admin.POST(pathLocations, s.createLocation)
	admin.GET(pathLocationsWithID, s.getLocation)
	admin.PUT(pathLocationsWithID, s.updateLocation)
	admin.DELETE(pathLocationsWithID, s.deleteLocation)

	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") && !strings.Contains(c.Request.RequestURI, ".") {
			file, _ := s.staticFiles.ReadFile("assets/index.html")
			c.Data(
				http.StatusOK,
				"text/html; charset=utf-8",
				file,
			)
		}
	})

	return r, nil
}

func (s *Server) registerCorsMiddleware(r *gin.Engine) {
	if len(s.cfg.App.CorsAllowOrigins) > 0 {
		slog.Info("CORS middleware enabled", "origins", s.cfg.App.CorsAllowOrigins)
		r.Use(s.createCorsMiddleware())
	} else {
		slog.Info("CORS middleware disabled (no origins configured)")
	}
}

func (s *Server) createCorsMiddleware() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = s.cfg.App.CorsAllowOrigins
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization", "Credentials", "Content-Type", "X-API-Key", "Accept")
	corsConfig.AddExposeHeaders("X-Total-Count", "Content-Disposition")

	return cors.New(corsConfig)
}
