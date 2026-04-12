package users

import (
	"context"
	"net/http"
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
	r.GET("/me", h.getMe)
}

func (h *Handler) getMe(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, _ := c.Get("user_id")
	
	user, err := h.service.GetProfile(ctx, userID.(string))
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "user not found", err)
		return
	}

	c.JSON(http.StatusOK, user)
}
