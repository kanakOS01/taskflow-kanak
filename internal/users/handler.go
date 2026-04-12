package users

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
	r.GET("/me", h.getMe)
}

func (h *Handler) getMe(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	user, err := h.service.GetProfile(userID.(string))
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "user not found")
		return
	}

	c.JSON(http.StatusOK, user)
}
