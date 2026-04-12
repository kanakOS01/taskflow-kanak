package utils

import "github.com/gin-gonic/gin"

type APIError struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func SendError(c *gin.Context, status int, msg string, err error) {
	if err != nil {
		c.Error(err)
	}
	c.AbortWithStatusJSON(status, APIError{Error: msg})
}

func SendValidationError(c *gin.Context, fields map[string]string, err error) {
	if err != nil {
		c.Error(err)
	}
	c.AbortWithStatusJSON(400, APIError{
		Error:  "validation failed",
		Fields: fields,
	})
}
