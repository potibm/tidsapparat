package hub

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/potibm/tidsapparat/internal/app/config"
	"github.com/potibm/tidsapparat/internal/app/domain"
	"github.com/potibm/tidsapparat/internal/app/exporter"
	"github.com/potibm/tidsapparat/internal/app/repository"
	"github.com/potibm/tidsapparat/internal/app/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Mock Repository ---

type mockScheduleEntryRepo struct {
	listResult    []domain.ScheduleEntry
	listTotal     int64
	listErr       error
	getByIDResult *domain.ScheduleEntry
	getByIDErr    error
	saveErr       error
	deleteErr     error
	lastSaved     *domain.ScheduleEntry
	lastDeletedID int64
}

func (m *mockScheduleEntryRepo) List(
	ctx context.Context,
	params repository.ScheduleEntryListParams,
	filters repository.ScheduleEntryListFilters,
) ([]domain.ScheduleEntry, int64, error) {
	return m.listResult, m.listTotal, m.listErr
}

func (m *mockScheduleEntryRepo) GetByID(ctx context.Context, id int64) (*domain.ScheduleEntry, error) {
	return m.getByIDResult, m.getByIDErr
}

func (m *mockScheduleEntryRepo) Save(ctx context.Context, entry *domain.ScheduleEntry) error {
	m.lastSaved = entry

	return m.saveErr
}

func (m *mockScheduleEntryRepo) Delete(ctx context.Context, id int64) error {
	m.lastDeletedID = id

	return m.deleteErr
}

func (m *mockScheduleEntryRepo) GetAllPreloaded(ctx context.Context) (domain.TimeTable, error) {
	return nil, nil
}

func (m *mockScheduleEntryRepo) GetByCategoryID(ctx context.Context, categoryID int64) ([]domain.ScheduleEntry, error) {
	return nil, nil
}

func (m *mockScheduleEntryRepo) GetByLocationID(ctx context.Context, locationID int64) ([]domain.ScheduleEntry, error) {
	return nil, nil
}

func setupScheduleEntryServer(t *testing.T, repo *mockScheduleEntryRepo) *Server {
	t.Helper()

	expMgr := exporter.NewManager(repo, 0)
	eventHub := services.NewEventHub(expMgr, nil, repo)

	return &Server{
		scheduleEntryRepo: repo,
		eventHub:          eventHub,
		cfg:               config.Config{App: config.AppConfig{Environment: "test"}},
		logger:            slog.New(slog.DiscardHandler),
	}
}

func TestListScheduleEntries(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{
		listResult: []domain.ScheduleEntry{
			{ID: 1, Title: "Entry 1"},
			{ID: 2, Title: "Entry 2"},
		},
		listTotal: 2,
	}
	srv := setupScheduleEntryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/schedule-entries", http.NoBody)

	srv.listScheduleEntries(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "2", w.Header().Get("X-Total-Count"))

	var result []domain.ScheduleEntry
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Len(t, result, 2)
}

func TestListScheduleEntries_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{listErr: errors.New("db error")}
	srv := setupScheduleEntryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/schedule-entries", http.NoBody)

	srv.listScheduleEntries(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetScheduleEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{
		getByIDResult: &domain.ScheduleEntry{ID: 42, Title: "Found"},
	}
	srv := setupScheduleEntryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/schedule-entries/42", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "42"}}

	srv.getScheduleEntry(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var result domain.ScheduleEntry
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, int64(42), result.ID)
	assert.Equal(t, "Found", result.Title)
}

func TestGetScheduleEntry_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{}
	srv := setupScheduleEntryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/schedule-entries/abc", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	srv.getScheduleEntry(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetScheduleEntry_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{getByIDErr: errors.New("not found")}
	srv := setupScheduleEntryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/schedule-entries/99", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "99"}}

	srv.getScheduleEntry(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateScheduleEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{
		getByIDErr: errors.New("not found"),
	}
	srv := setupScheduleEntryServer(t, repo)

	entry := domain.ScheduleEntry{Title: "New Entry"}
	body, _ := json.Marshal(entry)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/schedule-entries", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.createScheduleEntry(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "New Entry", repo.lastSaved.Title)
}

func TestCreateScheduleEntry_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{}
	srv := setupScheduleEntryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/schedule-entries", bytes.NewReader([]byte("not json")))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.createScheduleEntry(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateScheduleEntry_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{saveErr: errors.New("insert failed")}
	srv := setupScheduleEntryServer(t, repo)

	entry := domain.ScheduleEntry{Title: "New Entry"}
	body, _ := json.Marshal(entry)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/schedule-entries", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.createScheduleEntry(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateScheduleEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{
		getByIDErr: errors.New("not found"),
	}
	srv := setupScheduleEntryServer(t, repo)

	entry := domain.ScheduleEntry{Title: "Updated Entry"}
	body, _ := json.Marshal(entry)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/schedule-entries/5", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	srv.updateScheduleEntry(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(5), repo.lastSaved.ID)
	assert.Equal(t, "Updated Entry", repo.lastSaved.Title)
}

func TestUpdateScheduleEntry_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{}
	srv := setupScheduleEntryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/schedule-entries/abc", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	srv.updateScheduleEntry(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateScheduleEntry_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{saveErr: errors.New("update failed")}
	srv := setupScheduleEntryServer(t, repo)

	entry := domain.ScheduleEntry{Title: "Updated Entry"}
	body, _ := json.Marshal(entry)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/schedule-entries/5", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	srv.updateScheduleEntry(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteScheduleEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{}
	srv := setupScheduleEntryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/schedule-entries/7", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "7"}}

	srv.deleteScheduleEntry(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(7), repo.lastDeletedID)
}

func TestDeleteScheduleEntry_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{}
	srv := setupScheduleEntryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/schedule-entries/abc", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	srv.deleteScheduleEntry(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteScheduleEntry_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockScheduleEntryRepo{deleteErr: errors.New("delete failed")}
	srv := setupScheduleEntryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/schedule-entries/7", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "7"}}

	srv.deleteScheduleEntry(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestParseScheduleEntryListParams_Defaults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/schedule-entries", http.NoBody)

	params := parseScheduleEntryListParams(c)

	assert.Equal(t, 0, params.Offset)
	assert.Equal(t, 20, params.Limit)
	assert.Equal(t, "id", params.Sort)
	assert.Equal(t, "DESC", params.Order)
}

func TestParseScheduleEntryListParams_Custom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(
		http.MethodGet,
		"/schedule-entries?_start=10&_end=30&_sort=title&_order=ASC",
		http.NoBody,
	)

	params := parseScheduleEntryListParams(c)

	assert.Equal(t, 10, params.Offset)
	assert.Equal(t, 20, params.Limit)
	assert.Equal(t, "title", params.Sort)
	assert.Equal(t, "ASC", params.Order)
}

func TestParseScheduleEntryFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("all filters present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(
			http.MethodGet,
			"/schedule-entries?q=test&category_id=1&location_id=2&id=3&hide_past=true&hide_hidden=true",
			http.NoBody,
		)

		filters := parseScheduleEntryFilters(c)

		require.NotNil(t, filters.Query)
		assert.Equal(t, "test", *filters.Query)
		require.NotNil(t, filters.CategoryID)
		assert.Equal(t, int64(1), *filters.CategoryID)
		require.NotNil(t, filters.LocationID)
		assert.Equal(t, int64(2), *filters.LocationID)
		require.NotNil(t, filters.ID)
		assert.Equal(t, int64(3), *filters.ID)
		assert.True(t, filters.HidePast)
		assert.True(t, filters.HideHidden)
	})

	t.Run("no filters", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/schedule-entries", http.NoBody)

		filters := parseScheduleEntryFilters(c)

		assert.Nil(t, filters.Query)
		assert.Nil(t, filters.CategoryID)
		assert.Nil(t, filters.LocationID)
		assert.Nil(t, filters.ID)
		assert.False(t, filters.HidePast)
		assert.False(t, filters.HideHidden)
	})

	t.Run("invalid numeric ids are ignored", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(
			http.MethodGet,
			"/schedule-entries?category_id=abc&location_id=def&id=ghi",
			http.NoBody,
		)

		filters := parseScheduleEntryFilters(c)

		assert.Nil(t, filters.CategoryID)
		assert.Nil(t, filters.LocationID)
		assert.Nil(t, filters.ID)
	})
}
