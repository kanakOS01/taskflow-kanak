package tasks

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

func (h *Handler) RegisterRoutes(api *gin.RouterGroup) {
	// Nested under projects
	api.GET("/projects/:id/tasks", h.listByProject)
	api.POST("/projects/:id/tasks", h.create)
	
	// Direct task routes
	api.PATCH("/tasks/:id", h.update)
	api.DELETE("/tasks/:id", h.delete)
}

func (h *Handler) listByProject(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	projectID := c.Param("id")
	status := c.Query("status")
	assignee := c.Query("assignee")

	page := 1
	limit := 20

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	resp, err := h.service.List(ctx, projectID, status, assignee, page, limit)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")
	projectID := c.Param("id")

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, map[string]string{"body": err.Error()})
		return
	}

	task, err := h.service.Create(ctx, projectID, userID.(string), req)
	if err != nil {
		if err == ErrNotFound {
			utils.SendError(c, http.StatusNotFound, "project not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to create task")
		return
	}

	c.JSON(http.StatusCreated, TaskEnvelope{Task: task})
}

func (h *Handler) update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, map[string]string{"body": err.Error()})
		return
	}

	task, err := h.service.Update(ctx, id, userID.(string), req)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "unauthorized action")
			return
		}
		if err == ErrNotFound {
			utils.SendError(c, http.StatusNotFound, "task not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to update task")
		return
	}

	c.JSON(http.StatusOK, TaskEnvelope{Task: task})
}

func (h *Handler) delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")
	id := c.Param("id")

	err := h.service.Delete(ctx, id, userID.(string))
	if err != nil {
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "unauthorized action")
			return
		}
		if err == ErrNotFound {
			utils.SendError(c, http.StatusNotFound, "task not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to delete task")
		return
	}

	c.Status(http.StatusNoContent)
}
