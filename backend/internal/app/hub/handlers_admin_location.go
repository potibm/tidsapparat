package hub

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/potibm/billedapparat/internal/app/repository"
)

func (s *Server) listLocations(c *gin.Context) {
	params := parseLocationListParams(c)

	filters := parseLocationFilters(c)

	locations, total, err := s.locationRepo.List(c.Request.Context(), params, filters)
	if err != nil {
		respondWithInternalServerProblem(c, "Failed to list locations: "+err.Error())

		return
	}

	c.Header("X-Total-Count", strconv.FormatInt(total, 10))

	c.JSON(http.StatusOK, locations)
}

func (s *Server) getLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	location, err := s.locationRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		respondWithNotFoundProblem(c, "Location with ID "+strconv.FormatInt(id, 10)+" not found")

		return
	}

	c.JSON(http.StatusOK, location)
}

func (s *Server) createLocation(c *gin.Context) {
	var location domain.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		respondWithFailedToParsePayloadProblem(c, err)

		return
	}

	if err := s.locationRepo.Create(c.Request.Context(), &location); err != nil {
		slog.Error("Create Location: Failed to create location", "error", err)
		respondWithInternalServerProblem(c, "Failed to create location: "+err.Error())

		return
	}

	slog.Info("Create Location: Successfully created location", "id", location.ID)

	c.JSON(http.StatusCreated, location)
}

func (s *Server) updateLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	var location domain.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		respondWithFailedToParsePayloadProblem(c, err)

		return
	}

	location.ID = id

	if err := s.locationRepo.Update(c.Request.Context(), &location); err != nil {
		slog.Error("Update Location: Failed to update location", "id", id, "error", err)
		respondWithInternalServerProblem(c, "Failed to update location: "+err.Error())

		return
	}

	s.eventHub.SyncLocationUpdate(c.Request.Context(), location.ID)

	slog.Info("Update Location: Successfully updated location", "id", id)

	c.JSON(http.StatusOK, location)
}

func (s *Server) deleteLocation(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	if err := s.locationRepo.Delete(c.Request.Context(), id); err != nil {
		respondWithInternalServerProblem(c, "Failed to delete location: "+err.Error())

		return
	}

	s.eventHub.SyncLocationUpdate(c.Request.Context(), id)

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func parseLocationListParams(c *gin.Context) repository.LocationListParams {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "20"))

	return repository.LocationListParams{
		Offset: start,
		Limit:  end - start,
		Sort:   c.DefaultQuery("_sort", "id"),
		Order:  c.DefaultQuery("_order", "DESC"),
	}
}

func parseLocationFilters(c *gin.Context) repository.LocationListFilters {
	filters := repository.LocationListFilters{}

	if q := c.Query("q"); q != "" {
		filters.Query = &q
	}

	return filters
}
