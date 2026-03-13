package trello

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ClientOptions holds configuration for the API client.
type ClientOptions struct {
	Timeout        time.Duration
	MaxRetries     int
	RetryMutations bool
	Verbose        bool
}

// DefaultClientOptions returns sensible defaults.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Timeout:        15 * time.Second,
		MaxRetries:     3,
		RetryMutations: false,
		Verbose:        false,
	}
}

// Client is the Trello API client.
type Client struct {
	baseURL    string
	apiKey     string
	token      string
	httpClient *http.Client
	opts       ClientOptions
}

// NewClient creates a new Trello API client.
func NewClient(baseURL, apiKey, token string, opts ClientOptions) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		token:   token,
		httpClient: &http.Client{
			Timeout: opts.Timeout,
		},
		opts: opts,
	}
}

// Get performs an authenticated GET request and decodes the JSON response.
func (c *Client) Get(ctx context.Context, path string, params map[string]string, result any) error {
	return c.do(ctx, http.MethodGet, path, params, result)
}

// Post performs an authenticated POST request with params as query parameters.
// Trello API expects mutation data as query params, not JSON bodies.
func (c *Client) Post(ctx context.Context, path string, params map[string]string, result any) error {
	return c.do(ctx, http.MethodPost, path, params, result)
}

// Put performs an authenticated PUT request with params as query parameters.
func (c *Client) Put(ctx context.Context, path string, params map[string]string, result any) error {
	return c.do(ctx, http.MethodPut, path, params, result)
}

// Delete performs an authenticated DELETE request.
func (c *Client) Delete(ctx context.Context, path string, result any) error {
	return c.do(ctx, http.MethodDelete, path, nil, result)
}

// PostMultipart performs a multipart file upload POST (for attachments).
func (c *Client) PostMultipart(ctx context.Context, path string, params map[string]string, filePath string, result any) error {
	// Multipart upload implementation — handled in attachments.go
	return nil
}

// buildURL constructs the full URL with auth query params (key, token) and any
// additional request params. Shared by both do() and postMultipartFile().
func (c *Client) buildURL(path string, params map[string]string) string {
	u, _ := url.Parse(c.baseURL + path)
	q := u.Query()
	q.Set("key", c.apiKey)
	q.Set("token", c.token)
	for k, v := range params {
		if strings.Contains(v, ",") {
			for _, part := range strings.Split(v, ",") {
				q.Add(k, part)
			}
			continue
		}
		q.Add(k, v)
	}
	encoded := q.Encode()
	if encoded != "" {
		encoded = "&" + encoded
	}
	u.RawQuery = encoded
	return u.String()
}

func (c *Client) do(ctx context.Context, method, path string, params map[string]string, result any) error {
	fullURL := c.buildURL(path, params)

	maxAttempts := 1 + c.opts.MaxRetries
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	if isMutation(method) && !c.opts.RetryMutations {
		maxAttempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			if err := waitForRetry(ctx, attempt-1); err != nil {
				return err
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		start := time.Now()
		if c.opts.Verbose {
			logURL := c.baseURL + path
			slog.Debug("trello request", "method", method, "url", logURL, "attempt", attempt+1)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		if c.opts.Verbose {
			slog.Debug("trello response", "status", resp.StatusCode, "duration", time.Since(start), "attempt", attempt+1)
		}

		if resp.StatusCode >= http.StatusBadRequest {
			lastErr = mapHTTPError(resp)
			resp.Body.Close()
			if shouldRetry(resp.StatusCode) && attempt < maxAttempts-1 {
				continue
			}
			return lastErr
		}

		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				resp.Body.Close()
				return fmt.Errorf("failed to decode response: %w", err)
			}
		}
		resp.Body.Close()
		return nil
	}

	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("request failed after %d attempts", maxAttempts)
}
