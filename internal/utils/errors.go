package utils

import "github.com/gin-gonic/gin"

type APIError struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func SendError(c *gin.Context, status int, msg string) {
	c.JSON(status, APIError{Error: msg})
}

func SendValidationError(c *gin.Context, fields map[string]string) {
	c.JSON(400, APIError{
		Error:  "validation failed",
		Fields: fields,
	})
}
