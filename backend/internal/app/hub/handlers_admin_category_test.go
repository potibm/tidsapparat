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

type mockCategoryRepo struct {
	listResult    []domain.Category
	listTotal     int64
	listErr       error
	getByIDResult *domain.Category
	getByIDErr    error
	createErr     error
	updateErr     error
	deleteErr     error
	lastCreated   *domain.Category
	lastUpdated   *domain.Category
	lastDeletedID int64
}

func (m *mockCategoryRepo) List(
	ctx context.Context,
	params repository.CategoryListParams,
	filters repository.CategoryListFilters,
) ([]domain.Category, int64, error) {
	return m.listResult, m.listTotal, m.listErr
}

func (m *mockCategoryRepo) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	return m.getByIDResult, m.getByIDErr
}

func (m *mockCategoryRepo) Create(ctx context.Context, category *domain.Category) error {
	m.lastCreated = category

	return m.createErr
}

func (m *mockCategoryRepo) Update(ctx context.Context, category *domain.Category) error {
	m.lastUpdated = category

	return m.updateErr
}

func (m *mockCategoryRepo) Delete(ctx context.Context, id int64) error {
	m.lastDeletedID = id

	return m.deleteErr
}

// --- Noop ScheduleSource for EventHub ---

type noopScheduleSource struct{}

func (n *noopScheduleSource) GetByID(ctx context.Context, id int64) (*domain.ScheduleEntry, error) {
	return nil, errors.New("not found")
}

func (n *noopScheduleSource) GetAllPreloaded(ctx context.Context) (domain.TimeTable, error) {
	return nil, nil
}

func (n *noopScheduleSource) GetByCategoryID(ctx context.Context, categoryID int64) ([]domain.ScheduleEntry, error) {
	return nil, nil
}

func (n *noopScheduleSource) GetByLocationID(ctx context.Context, locationID int64) ([]domain.ScheduleEntry, error) {
	return nil, nil
}

func setupCategoryServer(t *testing.T, repo *mockCategoryRepo) *Server {
	t.Helper()

	expMgr := exporter.NewManager(&noopScheduleSource{}, 0)
	eventHub := services.NewEventHub(expMgr, nil, &noopScheduleSource{})

	return &Server{
		categoryRepo: repo,
		eventHub:     eventHub,
		cfg:          config.Config{App: config.AppConfig{Environment: "test"}},
		logger:       slog.New(slog.DiscardHandler),
	}
}

func TestListCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{
		listResult: []domain.Category{
			{ID: 1, Name: "Music", Color: "#FF0000"},
			{ID: 2, Name: "Graphics", Color: "#00FF00"},
		},
		listTotal: 2,
	}
	srv := setupCategoryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/categories", http.NoBody)

	srv.listCategories(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "2", w.Header().Get("X-Total-Count"))

	var result []domain.Category
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Len(t, result, 2)
}

func TestListCategories_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{listErr: errors.New("db error")}
	srv := setupCategoryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/categories", http.NoBody)

	srv.listCategories(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{
		getByIDResult: &domain.Category{ID: 42, Name: "Demo"},
	}
	srv := setupCategoryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/categories/42", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "42"}}

	srv.getCategory(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var result domain.Category
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, int64(42), result.ID)
	assert.Equal(t, "Demo", result.Name)
}

func TestGetCategory_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{}
	srv := setupCategoryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/categories/abc", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	srv.getCategory(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetCategory_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{getByIDErr: errors.New("not found")}
	srv := setupCategoryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/categories/99", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "99"}}

	srv.getCategory(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{}
	srv := setupCategoryServer(t, repo)

	category := domain.Category{Name: "New Category", Color: "#0000FF"}
	body, _ := json.Marshal(category)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.createCategory(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "New Category", repo.lastCreated.Name)
}

func TestCreateCategory_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{}
	srv := setupCategoryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader([]byte("not json")))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.createCategory(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateCategory_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{createErr: errors.New("insert failed")}
	srv := setupCategoryServer(t, repo)

	category := domain.Category{Name: "New Category"}
	body, _ := json.Marshal(category)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	srv.createCategory(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{}
	srv := setupCategoryServer(t, repo)

	category := domain.Category{Name: "Updated Category"}
	body, _ := json.Marshal(category)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/categories/5", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	srv.updateCategory(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(5), repo.lastUpdated.ID)
	assert.Equal(t, "Updated Category", repo.lastUpdated.Name)
}

func TestUpdateCategory_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{}
	srv := setupCategoryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/categories/abc", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	srv.updateCategory(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateCategory_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{updateErr: errors.New("update failed")}
	srv := setupCategoryServer(t, repo)

	category := domain.Category{Name: "Updated Category"}
	body, _ := json.Marshal(category)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/categories/5", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	srv.updateCategory(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{}
	srv := setupCategoryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/categories/7", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "7"}}

	srv.deleteCategory(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(7), repo.lastDeletedID)
}

func TestDeleteCategory_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{}
	srv := setupCategoryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/categories/abc", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	srv.deleteCategory(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteCategory_RepoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &mockCategoryRepo{deleteErr: errors.New("delete failed")}
	srv := setupCategoryServer(t, repo)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/categories/7", http.NoBody)
	c.Params = gin.Params{{Key: "id", Value: "7"}}

	srv.deleteCategory(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestParseCategoryListParams_Defaults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/categories", http.NoBody)

	params := parseCategoryListParams(c)

	assert.Equal(t, 0, params.Offset)
	assert.Equal(t, 20, params.Limit)
	assert.Equal(t, "id", params.Sort)
	assert.Equal(t, "DESC", params.Order)
}

func TestParseCategoryListParams_Custom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/categories?_start=5&_end=15&_sort=name&_order=ASC", http.NoBody)

	params := parseCategoryListParams(c)

	assert.Equal(t, 5, params.Offset)
	assert.Equal(t, 10, params.Limit)
	assert.Equal(t, "name", params.Sort)
	assert.Equal(t, "ASC", params.Order)
}

func TestParseCategoryFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("query filter present", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/categories?q=music", http.NoBody)

		filters := parseCategoryFilters(c)

		require.NotNil(t, filters.Query)
		assert.Equal(t, "music", *filters.Query)
	})

	t.Run("no filters", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest(http.MethodGet, "/categories", http.NoBody)

		filters := parseCategoryFilters(c)

		assert.Nil(t, filters.Query)
	})
}
