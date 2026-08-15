package platform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

var (
	resourceIDRegex      = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$`)
	clientReferenceRegex = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{2,127}$`)
	serviceCodeRegex     = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	providerIDRegex      = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)
	ianaTimeZoneRegex    = regexp.MustCompile(`^[A-Za-z_]+(/[A-Za-z0-9._+-]+)+$`)
	currencyRegex        = regexp.MustCompile(`^[A-Z]{3}$`)
)

const DefaultMaxBodyBytes = 1048576 // 1MB

var ErrInvalidJSON = errors.New("invalid JSON")

// DecodeStrictJSON decodes JSON strictly: rejecting unknown fields, bounding size, and checking for trailing data.
func DecodeStrictJSON[T any](bodyBytes []byte, maxBytes int64) (T, []ErrorDetail, error) {
	var zero T
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBodyBytes
	}
	if int64(len(bodyBytes)) > maxBytes {
		return zero, []ErrorDetail{{Reason: fmt.Sprintf("request body exceeds maximum limit of %d bytes", maxBytes)}}, ErrInvalidJSON
	}
	if len(bytes.TrimSpace(bodyBytes)) == 0 {
		return zero, []ErrorDetail{{Reason: "request body cannot be empty"}}, ErrInvalidJSON
	}

	dec := json.NewDecoder(bytes.NewReader(bodyBytes))
	dec.DisallowUnknownFields()

	var target T
	if err := dec.Decode(&target); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxErr):
			return zero, []ErrorDetail{{Reason: fmt.Sprintf("syntax error at offset %d", syntaxErr.Offset)}}, ErrInvalidJSON
		case errors.As(err, &unmarshalTypeErr):
			field := unmarshalTypeErr.Field
			if field == "" {
				field = "body"
			}
			return zero, []ErrorDetail{{Field: field, Reason: fmt.Sprintf("expected %s but received %s", unmarshalTypeErr.Type, unmarshalTypeErr.Value)}}, ErrInvalidJSON
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return zero, []ErrorDetail{{Field: fieldName, Reason: "unknown field is not allowed"}}, ErrInvalidJSON
		default:
			return zero, []ErrorDetail{{Reason: err.Error()}}, ErrInvalidJSON
		}
	}

	var extra json.RawMessage
	if err := dec.Decode(&extra); err != io.EOF {
		return zero, []ErrorDetail{{Reason: "unexpected trailing data after valid JSON"}}, ErrInvalidJSON
	}

	return target, nil, nil
}

// ValidateRFC3339UTC verifies that the string is a valid RFC3339 timestamp ending in literal 'Z'.
func ValidateRFC3339UTC(value string) bool {
	if !strings.HasSuffix(value, "Z") {
		return false
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t, err = time.Parse(time.RFC3339Nano, value)
	}
	return err == nil && t.Location() == time.UTC
}

// ValidateResourceId checks the resource ID pattern.
func ValidateResourceId(id string) bool {
	if len(id) < 3 || len(id) > 128 {
		return false
	}
	return resourceIDRegex.MatchString(id)
}

// ValidateClientReference checks the client reference pattern.
func ValidateClientReference(ref string) bool {
	if len(ref) < 3 || len(ref) > 128 {
		return false
	}
	return clientReferenceRegex.MatchString(ref)
}

// ValidateServiceCode checks the service code pattern.
func ValidateServiceCode(code string) bool {
	if len(code) < 1 || len(code) > 64 {
		return false
	}
	return serviceCodeRegex.MatchString(code)
}

// ValidateProviderId checks the provider ID pattern.
func ValidateProviderId(id string) bool {
	if len(id) < 3 || len(id) > 64 {
		return false
	}
	return providerIDRegex.MatchString(id)
}

// ValidateIanaTimezone checks the IANA timezone pattern.
func ValidateIanaTimezone(tz string) bool {
	if len(tz) < 3 || len(tz) > 64 {
		return false
	}
	return ianaTimeZoneRegex.MatchString(tz)
}

// ValidateMoney verifies that Money contains non-negative minor units and a valid 3-letter uppercase ISO currency.
func ValidateMoney(m Money, fieldPrefix string) []ErrorDetail {
	var details []ErrorDetail
	if m.AmountMinor < 0 {
		details = append(details, ErrorDetail{
			Field:  joinField(fieldPrefix, "amount_minor"),
			Reason: "must be non-negative",
		})
	}
	if !currencyRegex.MatchString(m.Currency) {
		details = append(details, ErrorDetail{
			Field:  joinField(fieldPrefix, "currency"),
			Reason: "must be a 3-letter uppercase ISO 4217 currency code",
		})
	}
	return details
}

// ValidateAccessibility checks bounds and uniqueness of an accessibility string array.
func ValidateAccessibility(list []string, field string) []ErrorDetail {
	var details []ErrorDetail
	if len(list) > 10 {
		details = append(details, ErrorDetail{
			Field:  field,
			Reason: "must not contain more than 10 items",
		})
	}
	seen := make(map[string]bool, len(list))
	for idx, item := range list {
		if len(item) < 1 || len(item) > 64 {
			details = append(details, ErrorDetail{
				Field:  fmt.Sprintf("%s[%d]", field, idx),
				Reason: "length must be between 1 and 64 characters",
			})
		}
		if seen[item] {
			details = append(details, ErrorDetail{
				Field:  fmt.Sprintf("%s[%d]", field, idx),
				Reason: fmt.Sprintf("duplicate item %q is not allowed", item),
			})
		}
		seen[item] = true
	}
	return details
}

func joinField(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}
