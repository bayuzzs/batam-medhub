package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const maxRequestBodyBytes = 64 * 1024 // 64 KB

// bindStrictJSON parses JSON with DisallowUnknownFields, single-value validation, and body-size limits.
func bindStrictJSON(c *gin.Context, dst any) error {
	if c.Request.Body == nil {
		return errors.New("request body is empty")
	}

	limitedReader := http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestBodyBytes)
	dec := json.NewDecoder(limitedReader)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	// Enforce exactly one JSON value (no trailing JSON data)
	var trailing any
	if err := dec.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("request body contains unexpected trailing JSON values")
		}
		return err
	}

	return nil
}
