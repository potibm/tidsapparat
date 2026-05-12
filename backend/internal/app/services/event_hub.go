package services

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/potibm/protokolapparat/pkg/common"
	"github.com/potibm/protokolapparat/pkg/schedule"
	"github.com/potibm/tidsapparat/internal/app/domain"
	"github.com/potibm/tidsapparat/internal/app/exporter"
	"github.com/redis/go-redis/v9"
)

const streamName = "party:schedule:events"

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

func (h *EventHub) PublishCreate(ctx context.Context, entryID int64) {
	entry, err := h.getProtocolEntry(ctx, entryID)
	if err == nil {
		h.send(ctx, common.NewCreateEvent(entry))
	}
}

func (h *EventHub) PublishUpdate(ctx context.Context, entryID int64) {
	entry, err := h.getProtocolEntry(ctx, entryID)
	if err == nil {
		h.send(ctx, common.NewUpdateEvent(entry))
	}
}

func (h *EventHub) PublishDelete(ctx context.Context, entryID int64) {
	h.send(ctx, common.NewDeleteEvent(schedule.Entry{ID: entryID}))
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

	mappedEntries := mapToTimeTablePayload(timetable)
	syncEvent := common.NewSyncEvent(mappedEntries)

	h.sendToStream(ctx, mappedEntries)

	h.logger.Info("Sent full state sync event to Redis", "count", len(syncEvent.Payload))
}

func (h *EventHub) SyncCategoryUpdate(ctx context.Context, catID int64) {
	entries, _ := h.repo.GetByCategoryID(ctx, catID)

	for _, entry := range entries {
		h.PublishUpdate(ctx, entry.ID)
	}

	h.exporter.Ping()
}

func (h *EventHub) SyncLocationUpdate(ctx context.Context, locationID int64) {
	entries, _ := h.repo.GetByLocationID(ctx, locationID)

	for _, entry := range entries {
		h.PublishUpdate(ctx, entry.ID)
	}

	h.exporter.Ping()
}

func (h *EventHub) getProtocolEntry(ctx context.Context, entryID int64) (schedule.Entry, error) {
	dbEntry, err := h.repo.GetByID(ctx, entryID)
	if err != nil {
		h.logger.Error("Failed to fetch schedule entry", "id", entryID, "error", err)

		return schedule.Entry{}, err
	}

	return mapToEventPayload(dbEntry), nil
}

func (h *EventHub) send(ctx context.Context, event common.Event[schedule.Entry]) {
	if h.redis == nil {
		return
	}

	if err := event.Validate(); err != nil {
		h.logger.Error("Tried to publish invalid event", "error", err, "action", event.Action)

		return
	}

	h.sendToStream(ctx, event)
	h.exporter.Ping()
}

func (h *EventHub) sendToStream(ctx context.Context, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Failed to marshal data for Redis", "error", err)

		return
	}

	err = h.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{"data": jsonData},
	}).Err()
	if err != nil {
		h.logger.Error("Redis XADD error", "error", err)
	}
}
