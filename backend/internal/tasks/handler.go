package tasks

import (
	"net/http"

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
	projectID := c.Param("id")
	status := c.Query("status")
	assignee := c.Query("assignee")

	tasks, err := h.service.List(projectID, status, assignee)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	projectID := c.Param("id")

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, map[string]string{"body": err.Error()})
		return
	}

	task, err := h.service.Create(projectID, userID.(string), req)
	if err != nil {
		if err == ErrNotFound {
			utils.SendError(c, http.StatusNotFound, "project not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to create task")
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *Handler) update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, map[string]string{"body": err.Error()})
		return
	}

	task, err := h.service.Update(id, userID.(string), req)
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

	c.JSON(http.StatusOK, task)
}

func (h *Handler) delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	err := h.service.Delete(id, userID.(string))
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

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
