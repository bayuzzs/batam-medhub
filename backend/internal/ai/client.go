package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrAIUnconfigured = errors.New("cloudflare workers ai is not configured")
	ErrAIRequestError = errors.New("cloudflare workers ai request failed")
)

// ChatMessage represents a single message in the chat prompt.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ModelClient defines the interface for calling an AI model.
type ModelClient interface {
	Infer(ctx context.Context, messages []ChatMessage) (string, error)
	IsConfigured() bool
}

// Client implements ModelClient using the Cloudflare Workers AI REST API.
type Client struct {
	baseURL    string
	accountID  string
	apiToken   string
	model      string
	httpClient *http.Client
}

// NewClient constructs a new Cloudflare Workers AI client.
func NewClient(baseURL, accountID, apiToken, model string, timeout time.Duration) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.cloudflare.com/client/v4"
	}
	if model == "" {
		model = "@cf/meta/llama-3.3-70b-instruct-fp8-fast"
	}
	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	return &Client{
		baseURL:   baseURL,
		accountID: strings.TrimSpace(accountID),
		apiToken:  strings.TrimSpace(apiToken),
		model:     strings.TrimSpace(model),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// IsConfigured returns true if enough configuration is present to call the
// model. OpenAI-compatible base URLs (containing "/ai/v1") carry the account in
// the URL itself, so only an API token is required there; the native
// Cloudflare REST API additionally needs an account ID.
func (c *Client) IsConfigured() bool {
	if c == nil || c.apiToken == "" {
		return false
	}
	if c.openAICompatible() {
		return true
	}
	return c.accountID != ""
}

// openAICompatible reports whether the configured base URL points at the
// OpenAI-compatible gateway (https://.../ai/v1), which uses the standard
// /chat/completions request and response shapes instead of the Cloudflare
// native /accounts/{id}/ai/run/{model} endpoint.
func (c *Client) openAICompatible() bool {
	return strings.Contains(c.baseURL, "/ai/v1")
}

type workersAIRequest struct {
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

type workersAIResponse struct {
	Success  bool     `json:"success"`
	Errors   []any    `json:"errors"`
	Messages []any    `json:"messages"`
	Result   any      `json:"result"`
}

// openAIRequest mirrors the OpenAI-compatible chat completions request.
type openAIRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

// openAIResponse mirrors the OpenAI-compatible chat completions response.
type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Infer sends a chat completion request to the configured AI endpoint.
func (c *Client) Infer(ctx context.Context, messages []ChatMessage) (string, error) {
	if !c.IsConfigured() {
		return "", ErrAIUnconfigured
	}

	if c.openAICompatible() {
		return c.inferOpenAI(ctx, messages)
	}
	return c.inferWorkersAI(ctx, messages)
}

// inferWorkersAI uses the native Cloudflare Workers AI REST API:
// POST {base}/accounts/{accountID}/ai/run/{model}.
func (c *Client) inferWorkersAI(ctx context.Context, messages []ChatMessage) (string, error) {
	reqURL := fmt.Sprintf("%s/accounts/%s/ai/run/%s", c.baseURL, c.accountID, c.model)

	reqPayload := workersAIRequest{
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 0.1,
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("marshal ai request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create ai http request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrAIRequestError, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read ai response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("%w (status %d): %s", ErrAIRequestError, resp.StatusCode, string(respBody))
	}

	var aiResp workersAIResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		return "", fmt.Errorf("decode workers ai response envelope: %w", err)
	}

	if !aiResp.Success && len(aiResp.Errors) > 0 {
		errBytes, _ := json.Marshal(aiResp.Errors)
		return "", fmt.Errorf("%w: cloudflare error: %s", ErrAIRequestError, string(errBytes))
	}

	// Result can be map with "response" key or string directly
	switch res := aiResp.Result.(type) {
	case string:
		return res, nil
	case map[string]any:
		if textVal, ok := res["response"].(string); ok {
			return textVal, nil
		}
		if textVal, ok := res["text"].(string); ok {
			return textVal, nil
		}
		// If whole result object is JSON
		rawJSON, _ := json.Marshal(res)
		return string(rawJSON), nil
	default:
		rawJSON, _ := json.Marshal(res)
		return string(rawJSON), nil
	}
}

// inferOpenAI uses the OpenAI-compatible chat completions endpoint:
// POST {base}/chat/completions.
func (c *Client) inferOpenAI(ctx context.Context, messages []ChatMessage) (string, error) {
	reqURL := c.baseURL + "/chat/completions"

	reqPayload := openAIRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 0.1,
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("marshal ai request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("create ai http request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiToken)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrAIRequestError, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read ai response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("%w (status %d): %s", ErrAIRequestError, resp.StatusCode, string(respBody))
	}

	var aiResp openAIResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		return "", fmt.Errorf("decode openai response envelope: %w", err)
	}

	if aiResp.Error != nil {
		return "", fmt.Errorf("%w: %s", ErrAIRequestError, aiResp.Error.Message)
	}

	var builder strings.Builder
	for _, choice := range aiResp.Choices {
		builder.WriteString(choice.Message.Content)
	}
	return builder.String(), nil
}
