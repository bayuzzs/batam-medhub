package platform

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"regexp"

	"github.com/gin-gonic/gin"
)

var idempotencyKeyRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{7,127}$`)

// ValidateIdempotencyKey validates the Idempotency-Key header format.
func ValidateIdempotencyKey(key string) bool {
	return idempotencyKeyRegex.MatchString(key)
}

// ComputeFingerprint computes a SHA-256 hash of the method, path, and request body.
func ComputeFingerprint(method, path string, body []byte) string {
	hasher := sha256.New()
	hasher.Write([]byte(method))
	hasher.Write([]byte(" "))
	hasher.Write([]byte(path))
	hasher.Write([]byte("\n"))

	// Canonicalize JSON if valid JSON to ignore formatting/whitespace differences
	var obj any
	if len(body) > 0 && json.Unmarshal(body, &obj) == nil {
		if canonical, err := json.Marshal(obj); err == nil {
			hasher.Write(canonical)
			return hex.EncodeToString(hasher.Sum(nil))
		}
	}

	hasher.Write(body)
	return hex.EncodeToString(hasher.Sum(nil))
}

// ReadAndRestoreBody reads the request body and restores it so subsequent handlers can read it.
func ReadAndRestoreBody(c *gin.Context) ([]byte, error) {
	if c.Request.Body == nil {
		return []byte{}, nil
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}
