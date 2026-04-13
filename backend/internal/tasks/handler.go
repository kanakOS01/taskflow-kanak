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
	api.GET("/projects/:id/tasks", h.listByProject)
	api.POST("/projects/:id/tasks", h.create)

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
		if err == ErrProjectNotFound {
			utils.SendError(c, http.StatusNotFound, "project not found", err)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch tasks", err)
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
		utils.SendValidationError(c, map[string]string{"body": err.Error()}, err)
		return
	}

	task, err := h.service.Create(ctx, projectID, userID.(string), req)
	if err != nil {
		switch err {
		case ErrProjectNotFound:
			utils.SendError(c, http.StatusNotFound, "project not found", err)
		case ErrForbidden:
			utils.SendError(c, http.StatusForbidden, "only the project owner can create tasks", err)
		case ErrPastDueDate:
			utils.SendError(c, http.StatusUnprocessableEntity, err.Error(), err)
		case ErrInvalidAssignee:
			utils.SendError(c, http.StatusUnprocessableEntity, "assignee_id does not refer to a valid user", err)
		default:
			utils.SendError(c, http.StatusInternalServerError, "failed to create task", err)
		}
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *Handler) update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req UpdateTaskRequest
	if err := utils.BindStrict(c, &req); err != nil {
		utils.SendValidationError(c, map[string]string{"body": err.Error()}, err)
		return
	}

	task, err := h.service.Update(ctx, id, userID.(string), req)
	if err != nil {
		switch err {
		case ErrForbidden:
			utils.SendError(c, http.StatusForbidden, "unauthorized action", err)
		case ErrNotFound:
			utils.SendError(c, http.StatusNotFound, "task not found", err)
		case ErrProjectNotFound:
			utils.SendError(c, http.StatusNotFound, "project not found", err)
		case ErrPastDueDate:
			utils.SendError(c, http.StatusUnprocessableEntity, err.Error(), err)
		case ErrInvalidAssignee:
			utils.SendError(c, http.StatusUnprocessableEntity, "assignee_id does not refer to a valid user", err)
		default:
			utils.SendError(c, http.StatusInternalServerError, "failed to update task", err)
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")
	id := c.Param("id")

	err := h.service.Delete(ctx, id, userID.(string))
	if err != nil {
		switch err {
		case ErrForbidden:
			utils.SendError(c, http.StatusForbidden, "unauthorized action", err)
		case ErrNotFound:
			utils.SendError(c, http.StatusNotFound, "task not found", err)
		case ErrProjectNotFound:
			utils.SendError(c, http.StatusNotFound, "project not found", err)
		default:
			utils.SendError(c, http.StatusInternalServerError, "failed to delete task", err)
		}
		return
	}

	c.Status(http.StatusNoContent)
}
