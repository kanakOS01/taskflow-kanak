package utils

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// BindStrict decodes the request body into dst using DisallowUnknownFields,
// so any unexpected JSON keys are rejected with a clear error.
func BindStrict(c *gin.Context, dst interface{}) error {
	body, err := c.GetRawData()
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid request body: %w", err)
	}
	return nil
}
