package projects

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

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.list)
	r.POST("", h.create)
	r.GET("/:id", h.getByID)
	r.PATCH("/:id", h.update)
	r.DELETE("/:id", h.delete)
}

func (h *Handler) list(c *gin.Context) {
	userID, _ := c.Get("user_id")
	projects, err := h.service.List(userID.(string))
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch projects")
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *Handler) create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, map[string]string{"body": err.Error()})
		return
	}

	project, err := h.service.Create(userID.(string), req)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "failed to create project")
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *Handler) getByID(c *gin.Context) {
	id := c.Param("id")
	details, err := h.service.GetDetails(id)
	if err != nil {
		if err == ErrNotFound {
			utils.SendError(c, http.StatusNotFound, "project not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to fetch project")
		return
	}

	c.JSON(http.StatusOK, details)
}

func (h *Handler) update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, map[string]string{"body": err.Error()})
		return
	}

	project, err := h.service.Update(id, userID.(string), req)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.SendError(c, http.StatusForbidden, "unauthorized action")
			return
		}
		if err == ErrNotFound {
			utils.SendError(c, http.StatusNotFound, "project not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to update project")
		return
	}

	c.JSON(http.StatusOK, project)
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
			utils.SendError(c, http.StatusNotFound, "project not found")
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "failed to delete project")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
