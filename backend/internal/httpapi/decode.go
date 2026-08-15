package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxRequestBodyBytes = 64 * 1024 // 64 KB

var emailSchemaRegex = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)

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

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	if strings.ContainsAny(email, "<>\"(),:;[]\\ \t\r\n") {
		return false
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Name != "" || parsed.Address != email {
		return false
	}
	return emailSchemaRegex.MatchString(email)
}

func isValidPassword(password string) bool {
	length := len([]byte(password))
	return length >= 8 && length <= 72
}

func isValidRefreshToken(token string) bool {
	token = strings.TrimSpace(token)
	return len(token) >= 43 && len(token) <= 256
}
