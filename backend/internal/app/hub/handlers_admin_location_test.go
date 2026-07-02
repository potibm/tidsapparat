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

type mockLocationRepo struct {
	listResult    []domain.Location
	listTotal     int64
	listErr       error
	getByIDResult *domain.Location
	getByIDErr    error
	createErr     error
	updateErr     error
	deleteErr     error
	lastCreated   *domain.Location
	lastUpdated   *domain.Location
	lastDeletedID int64
}

func (m *mockLocationRepo) List(
	ctx context.Context,
	params repository.LocationListParams,
	filters repository.LocationListFilters,
) ([]domain.Location, int64, error) {
	return m.listResult, m.listTotal, m.listErr
}

func (m *mockLocationRepo) GetByID(ctx context.Context, id int64) (*domain.Location, error) {
	return m.getByIDResult, m.getByIDErr
}

func (m *mockLocationRepo) Create(ctx context.Context, location *domain.Location) error {
	m.lastCreated = location

	return m.createErr
}

func (m *mockLocationRepo) Update(ctx context.Context, location *domain.Location) error {
	m.lastUpdated = location

	return m.updateErr
}

func (m *mockLocationRepo) Delete(ctx context.Context, id int64) error {
	m.lastDeletedID = id

	return m.deleteErr
}

func setupLocationServer(t *testing.T, repo *mockLocationRepo) *Server {
	t.Helper()

	expMgr := exporter.NewManager(&noopScheduleSource{}, 0)
	eventHub := services.NewEventHub(expMgr, nil, &noopScheduleSource{})

	return &Server{
		locationRepo: repo,
		eventHub:     eventHub,
		cfg:          config.Config{App: config.AppConfig{Environment: "test"}},
		logger:       slog.New(slog.DiscardHandler),
	}
}

func TestListLocations(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{
		listResult: []domain.Location{
			{ID: 1, Name: "Main Hall"},
			{ID: 2, Name: "Room A"},
		},
		listTotal: 2,
	}
	srv := setupLocationServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/locations", http.NoBody)

	srv.listLocations(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "2", w.Header().Get("X-Total-Count"))

	var result []domain.Location
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Len(t, result, 2)
}

func TestListLocations_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{listErr: errors.New("db error")}
	srv := setupLocationServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/locations", http.NoBody)

	srv.listLocations(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{
		getByIDResult: &domain.Location{ID: 42, Name: "Venue"},
	}
	srv := setupLocationServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/locations/42", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "42"}}

	srv.getLocation(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var result domain.Location
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, int64(42), result.ID)
	assert.Equal(t, "Venue", result.Name)
}

func TestGetLocation_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{}
	srv := setupLocationServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/locations/abc", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	srv.getLocation(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetLocation_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{getByIDErr: errors.New("not found")}
	srv := setupLocationServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/locations/99", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "99"}}

	srv.getLocation(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{}
	srv := setupLocationServer(t, repo)

	location := domain.Location{Name: "New Venue"}
	body, _ := json.Marshal(location)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/locations", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.createLocation(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "New Venue", repo.lastCreated.Name)
}

func TestCreateLocation_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{}
	srv := setupLocationServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/locations", bytes.NewReader([]byte("not json")))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.createLocation(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateLocation_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{createErr: errors.New("insert failed")}
	srv := setupLocationServer(t, repo)

	location := domain.Location{Name: "New Venue"}
	body, _ := json.Marshal(location)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/locations", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.createLocation(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{}
	srv := setupLocationServer(t, repo)

	location := domain.Location{Name: "Updated Venue"}
	body, _ := json.Marshal(location)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/locations/5", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	srv.updateLocation(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(5), repo.lastUpdated.ID)
	assert.Equal(t, "Updated Venue", repo.lastUpdated.Name)
}

func TestUpdateLocation_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{}
	srv := setupLocationServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/locations/abc", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	srv.updateLocation(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateLocation_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{updateErr: errors.New("update failed")}
	srv := setupLocationServer(t, repo)

	location := domain.Location{Name: "Updated Venue"}
	body, _ := json.Marshal(location)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/locations/5", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	srv.updateLocation(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteLocation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{}
	srv := setupLocationServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/locations/7", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "7"}}

	srv.deleteLocation(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(7), repo.lastDeletedID)
}

func TestDeleteLocation_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{}
	srv := setupLocationServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/locations/abc", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	srv.deleteLocation(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteLocation_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockLocationRepo{deleteErr: errors.New("delete failed")}
	srv := setupLocationServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/locations/7", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "7"}}

	srv.deleteLocation(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestParseLocationListParams_Defaults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/locations", http.NoBody)

	params := parseLocationListParams(c)

	assert.Equal(t, 0, params.Offset)
	assert.Equal(t, 20, params.Limit)
	assert.Equal(t, "id", params.Sort)
	assert.Equal(t, "DESC", params.Order)
}

func TestParseLocationListParams_Custom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/locations?_start=5&_end=15&_sort=name&_order=ASC", http.NoBody)

	params := parseLocationListParams(c)

	assert.Equal(t, 5, params.Offset)
	assert.Equal(t, 10, params.Limit)
	assert.Equal(t, "name", params.Sort)
	assert.Equal(t, "ASC", params.Order)
}

func TestParseLocationFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("query filter present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/locations?q=hall", http.NoBody)

		filters := parseLocationFilters(c)

		require.NotNil(t, filters.Query)
		assert.Equal(t, "hall", *filters.Query)
	})

	t.Run("no filters", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/locations", http.NoBody)

		filters := parseLocationFilters(c)

		assert.Nil(t, filters.Query)
	})
}
