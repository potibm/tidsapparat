package hub

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRespondWithProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test-path", http.NoBody)

	respondWithProblem(c, http.StatusBadRequest, "Bad Request", "Something went wrong")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))

	var result ProblemDetails
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, "about:blank", result.Type)
	assert.Equal(t, "Bad Request", result.Title)
	assert.Equal(t, http.StatusBadRequest, result.Status)
	assert.Equal(t, "Something went wrong", result.Detail)
	assert.Equal(t, "/test-path", result.Instance)
}

func TestRespondWithInvalidIDFormatProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/items/abc", http.NoBody)

	respondWithInvalidIDFormatProblem(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result ProblemDetails
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, "Bad Request", result.Title)
	assert.Equal(t, "Invalid ID format", result.Detail)
}

func TestRespondWithInternalServerProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/items", http.NoBody)

	respondWithInternalServerProblem(c, "Database connection lost")

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var result ProblemDetails
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, "Internal Server Error", result.Title)
	assert.Equal(t, "Database connection lost", result.Detail)
}

func TestRespondWithNotFoundProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/items/99", http.NoBody)

	respondWithNotFoundProblem(c, "Item 99 not found")

	assert.Equal(t, http.StatusNotFound, w.Code)

	var result ProblemDetails
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, "Not Found", result.Title)
	assert.Equal(t, "Item 99 not found", result.Detail)
}

func TestRespondWithFailedToParsePayloadProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/items", http.NoBody)

	respondWithFailedToParsePayloadProblem(c, assert.AnError)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result ProblemDetails
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, "Bad Request", result.Title)
	assert.Contains(t, result.Detail, "Failed to parse payload")
}

func TestRespondWithBadRequestProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/items", http.NoBody)

	respondWithBadRequestProblem(c, "Missing required field 'name'")

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var result ProblemDetails
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, "Bad Request", result.Title)
	assert.Equal(t, "Missing required field 'name'", result.Detail)
}

func TestRespondWithUnauthorizedProblem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/items", http.NoBody)

	respondWithUnauthorizedProblem(c, "Invalid token")

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var result ProblemDetails
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &result))
	assert.Equal(t, "Unauthorized", result.Title)
	assert.Equal(t, "Invalid token", result.Detail)
}
