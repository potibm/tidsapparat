package hub

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/potibm/billedapparat/internal/app/repository"
)

func (s *Server) listCategories(c *gin.Context) {
	params := parseCategoryListParams(c)

	filters := parseCategoryFilters(c)

	categories, total, err := s.categoryRepo.List(c.Request.Context(), params, filters)
	if err != nil {
		respondWithInternalServerProblem(c, "Failed to list categories: "+err.Error())

		return
	}

	c.Header("X-Total-Count", strconv.FormatInt(total, 10))

	c.JSON(http.StatusOK, categories)
}

func (s *Server) getCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	category, err := s.categoryRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		respondWithNotFoundProblem(c, "Category with ID "+strconv.FormatInt(id, 10)+" not found")

		return
	}

	c.JSON(http.StatusOK, category)
}

func (s *Server) createCategory(c *gin.Context) {
	var category domain.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		respondWithFailedToParsePayloadProblem(c, err)

		return
	}

	if err := s.categoryRepo.Create(c.Request.Context(), &category); err != nil {
		slog.Error("Create Category: Failed to create category", "error", err)
		respondWithInternalServerProblem(c, "Failed to create category: "+err.Error())

		return
	}

	slog.Info("Create Category: Successfully created category", "id", category.ID)

	c.JSON(http.StatusCreated, category)
}

func (s *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	var category domain.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		respondWithFailedToParsePayloadProblem(c, err)

		return
	}

	category.ID = id

	if err := s.categoryRepo.Update(c.Request.Context(), &category); err != nil {
		slog.Error("Update Category: Failed to update category", "id", id, "error", err)
		respondWithInternalServerProblem(c, "Failed to update category: "+err.Error())

		return
	}

	s.eventHub.SyncCategoryUpdate(c.Request.Context(), category.ID)

	slog.Info("Update Category: Successfully updated category", "id", id)

	c.JSON(http.StatusOK, category)
}

func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	if err := s.categoryRepo.Delete(c.Request.Context(), id); err != nil {
		respondWithInternalServerProblem(c, "Failed to delete category: "+err.Error())

		return
	}

	s.eventHub.SyncCategoryUpdate(c.Request.Context(), id)

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func parseCategoryListParams(c *gin.Context) repository.CategoryListParams {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "20"))

	return repository.CategoryListParams{
		Offset: start,
		Limit:  end - start,
		Sort:   c.DefaultQuery("_sort", "id"),
		Order:  c.DefaultQuery("_order", "DESC"),
	}
}

func parseCategoryFilters(c *gin.Context) repository.CategoryListFilters {
	filters := repository.CategoryListFilters{}

	if q := c.Query("q"); q != "" {
		filters.Query = &q
	}

	return filters
}
