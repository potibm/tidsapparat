package hub

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/potibm/billedapparat/internal/app/repository"
	"github.com/potibm/billedapparat/internal/app/services"
)

func (s *Server) listScheduleEntries(c *gin.Context) {
	params := parseScheduleEntryListParams(c)

	filters := parseScheduleEntryFilters(c)

	slides, total, err := s.scheduleEntryRepo.List(c.Request.Context(), params, filters)
	if err != nil {
		respondWithInternalServerProblem(c, "Failed to list slides: "+err.Error())

		return
	}

	c.Header("X-Total-Count", strconv.FormatInt(total, 10))

	c.JSON(http.StatusOK, slides)
}

func (s *Server) getScheduleEntry(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	slide, err := s.scheduleEntryRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		respondWithNotFoundProblem(c, "Schedule entry with ID "+strconv.FormatInt(id, 10)+" not found")

		return
	}

	c.JSON(http.StatusOK, slide)
}

func (s *Server) createScheduleEntry(c *gin.Context) {
	var scheduleEntry domain.ScheduleEntry
	if err := c.ShouldBindJSON(&scheduleEntry); err != nil {
		respondWithFailedToParsePayloadProblem(c, err)

		return
	}

	if err := s.scheduleEntryRepo.Save(c.Request.Context(), &scheduleEntry); err != nil {
		slog.Error("Create Schedule Entry: Failed to create slide", "error", err)
		respondWithInternalServerProblem(c, "Failed to create slide: "+err.Error())

		return
	} else {
		s.eventHub.Publish(c, scheduleEntry.ID, services.ActionCreate)

		slog.Info("Create Schedule Entry: Successfully created schedule entry", "id", scheduleEntry.ID)
	}

	c.JSON(http.StatusCreated, scheduleEntry)
}

func (s *Server) updateScheduleEntry(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	var scheduleEntry domain.ScheduleEntry
	if err := c.ShouldBindJSON(&scheduleEntry); err != nil {
		respondWithFailedToParsePayloadProblem(c, err)

		return
	}

	scheduleEntry.ID = id

	if err := s.scheduleEntryRepo.Save(c.Request.Context(), &scheduleEntry); err != nil {
		slog.Error("Update Schedule Entry: Failed to update schedule entry", "id", id, "error", err)
		respondWithInternalServerProblem(c, "Failed to update schedule entry: "+err.Error())

		return
	} else {
		s.eventHub.Publish(c, scheduleEntry.ID, services.ActionUpdate)

		slog.Info("Update Schedule Entry: Successfully updated schedule entry", "id", id)
	}

	c.JSON(http.StatusOK, scheduleEntry)
}

func (s *Server) deleteScheduleEntry(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondWithInvalidIDFormatProblem(c)

		return
	}

	if err := s.scheduleEntryRepo.Delete(c.Request.Context(), id); err != nil {
		respondWithInternalServerProblem(c, "Failed to delete schedule entry: "+err.Error())

		return
	}

	s.eventHub.Publish(c, id, services.ActionDelete)

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func parseScheduleEntryListParams(c *gin.Context) repository.ScheduleEntryListParams {
	start, _ := strconv.Atoi(c.DefaultQuery("_start", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("_end", "20"))

	return repository.ScheduleEntryListParams{
		Offset: start,
		Limit:  end - start,
		Sort:   c.DefaultQuery("_sort", "id"),
		Order:  c.DefaultQuery("_order", "DESC"),
	}
}

func parseScheduleEntryFilters(c *gin.Context) repository.ScheduleEntryListFilters {
	filters := repository.ScheduleEntryListFilters{}

	if q := c.Query("q"); q != "" {
		filters.Query = &q
	}

	if catIDStr := c.Query("category_id"); catIDStr != "" {
		if catID, err := strconv.ParseInt(catIDStr, 10, 64); err == nil {
			filters.CategoryID = &catID
		}
	}

	if locIDStr := c.Query("location_id"); locIDStr != "" {
		if locID, err := strconv.ParseInt(locIDStr, 10, 64); err == nil {
			filters.LocationID = &locID
		}
	}

	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			filters.ID = &id
		}
	}

	filters.HidePast = false
	if c.Query("hide_past") == "true" {
		filters.HidePast = true
	}

	return filters
}
