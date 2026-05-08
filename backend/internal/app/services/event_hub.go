package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/potibm/billedapparat/internal/app/exporter"
	"github.com/redis/go-redis/v9"
)

type ScheduleSource interface {
	GetByCategoryID(ctx context.Context, categoryID int64) ([]domain.ScheduleEntry, error)
	GetByLocationID(ctx context.Context, locationID int64) ([]domain.ScheduleEntry, error)
	GetAllPreloaded(ctx context.Context) (domain.TimeTable, error)
	GetByID(ctx context.Context, id int64) (*domain.ScheduleEntry, error)
}

type EventHub struct {
	exporter *exporter.Manager
	redis    *redis.Client
	repo     ScheduleSource
	logger   *slog.Logger
}

func NewEventHub(e *exporter.Manager, redisClient *redis.Client, repo ScheduleSource) *EventHub {
	logger := slog.Default().With("component", "EventHub")

	return &EventHub{
		exporter: e,
		redis:    redisClient,
		repo:     repo,
		logger:   logger,
	}
}

func (h *EventHub) Publish(ctx context.Context, entryID int64, action ActionType) {
	if h.redis != nil {
		var eventDTO ScheduleEventDTO

		if action == ActionDelete {
			eventDTO = ScheduleEventDTO{
				Action:    action,
				Timestamp: time.Now().Unix(),
				Payload:   ScheduleEntryDTO{ID: entryID},
			}
		} else {
			entry, err := h.repo.GetByID(ctx, entryID)
			if err != nil {
				h.logger.Error("Failed to fetch schedule entry", "id", entryID, "error", err)

				return
			}

			eventDTO = mapToEventDTO(entry, action)
		}

		h.sendToStream(ctx, eventDTO)
	}

	h.exporter.Ping()
}

func (h *EventHub) PublishFullSync(ctx context.Context) {
	if h.redis == nil {
		return
	}

	timetable, err := h.repo.GetAllPreloaded(ctx)
	if err != nil {
		h.logger.Error("Failed to fetch timetable for sync", "error", err)

		return
	}

	syncEvent := ScheduleSyncEventDTO{
		Action:    ActionSync,
		Timestamp: time.Now().Unix(),
		Payload:   mapToTimeTableDTO(timetable),
	}

	h.sendToStream(ctx, syncEvent)

	h.logger.Info("Sent full state sync event to Redis", "count", len(syncEvent.Payload))
}

func (h *EventHub) SyncCategoryUpdate(ctx context.Context, catID int64) {
	entries, _ := h.repo.GetByCategoryID(ctx, catID)

	for _, entry := range entries {
		h.Publish(ctx, entry.ID, "updated")
	}

	h.exporter.Ping()
}

func (h *EventHub) SyncLocationUpdate(ctx context.Context, locationID int64) {
	entries, _ := h.repo.GetByLocationID(ctx, locationID)

	for _, entry := range entries {
		h.Publish(ctx, entry.ID, "updated")
	}

	h.exporter.Ping()
}

func (h *EventHub) sendToStream(ctx context.Context, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Failed to marshal data for Redis", "error", err)

		return
	}

	err = h.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: "party:schedule:events",
		Values: map[string]interface{}{"data": jsonData},
	}).Err()
	if err != nil {
		h.logger.Error("Redis XADD error", "error", err)
	}
}
