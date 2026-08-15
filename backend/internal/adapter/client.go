package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"batam-medhub/internal/auth"
)

const maxResponseBodyBytes = 2 * 1024 * 1024 // 2MB

// Client is the base HTTP client for communicating with mock provider services.
type Client struct {
	baseURL        string
	integrationKey string
	httpClient     *http.Client
}

// NewClient constructs a new provider base client.
func NewClient(baseURL, integrationKey string, timeout time.Duration) *Client {
	return &Client{
		baseURL:        strings.TrimRight(baseURL, "/"),
		integrationKey: integrationKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Do executes an HTTP request against a provider service, attaching integration auth, request ID, and idempotency headers.
func (c *Client) Do(ctx context.Context, method, path, reqID, idemKey string, reqBody any, dest any) error {
	url := c.baseURL + path

	var bodyReader io.Reader
	if reqBody != nil {
		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("marshal provider request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create provider http request: %w", err)
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set Integration Secret
	req.Header.Set("X-Integration-Key", c.integrationKey)

	// Set Request ID
	if reqID == "" {
		reqID = "req-" + auth.NewUUID()
	}
	req.Header.Set("X-Request-ID", reqID)

	// Set Idempotency Key
	if idemKey != "" {
		req.Header.Set("Idempotency-Key", idemKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute provider request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyBytes))
	if err != nil {
		return fmt.Errorf("read provider response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errEnv ErrorEnvelope
		if err := json.Unmarshal(respBytes, &errEnv); err == nil && errEnv.Error.Code != "" {
			return &ProviderError{
				StatusCode: resp.StatusCode,
				Code:       errEnv.Error.Code,
				Message:    errEnv.Error.Message,
				Retryable:  errEnv.Error.Retryable,
				RequestID:  errEnv.Error.RequestID,
				Details:    errEnv.Error.Details,
			}
		}
		return &ProviderError{
			StatusCode: resp.StatusCode,
			Code:       fmt.Sprintf("HTTP_%d", resp.StatusCode),
			Message:    string(respBytes),
			Retryable:  resp.StatusCode >= 500,
			RequestID:  reqID,
			Details:    []ErrorDetail{},
		}
	}

	if dest != nil && len(respBytes) > 0 {
		decoder := json.NewDecoder(bytes.NewReader(respBytes))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(dest); err != nil {
			// Fallback to standard decode if provider has minor extra fields
			if err := json.Unmarshal(respBytes, dest); err != nil {
				return fmt.Errorf("decode provider response: %w", err)
			}
		}
	}

	return nil
}
