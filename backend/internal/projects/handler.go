package projects

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"taskflow/internal/utils"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.list)
	r.POST("", h.create)
	r.GET("/:id", h.getByID)
	r.GET("/:id/stats", h.stats)
	r.PATCH("/:id", h.update)
	r.DELETE("/:id", h.delete)
}

func (h *Handler) list(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			if v > 100 {
				v = 100
			}
			limit = v
		}
	}

	userID, _ := c.Get("user_id")
	resp, err := h.service.List(ctx, userID.(string), page, limit)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch projects")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, map[string]string{"body": err.Error()})
		return
	}

	project, err := h.service.Create(ctx, userID.(string), req)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to create project")
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *Handler) getByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	id := c.Param("id")
	details, err := h.service.GetDetails(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			utils.SendError(c, http.StatusNotFound, "not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch project")
		return
	}

	c.JSON(http.StatusOK, details)
}

func (h *Handler) stats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	id := c.Param("id")
	stats, err := h.service.GetStats(ctx, id)
	if err != nil {
		if err == ErrNotFound {
			utils.SendError(c, http.StatusNotFound, "not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch project stats")
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, map[string]string{"body": err.Error()})
		return
	}

	project, err := h.service.Update(ctx, id, userID.(string), req)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "forbidden")
			return
		}
		if err == ErrNotFound {
			utils.SendError(c, http.StatusNotFound, "not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to update project")
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *Handler) delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")
	id := c.Param("id")

	err := h.service.Delete(ctx, id, userID.(string))
	if err != nil {
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "forbidden")
			return
		}
		if err == ErrNotFound {
			utils.SendError(c, http.StatusNotFound, "not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to delete project")
		return
	}

	c.Status(http.StatusNoContent)
}
