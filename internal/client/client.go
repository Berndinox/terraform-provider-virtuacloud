package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://api.virtua.cloud/v1"
	defaultTimeout = 60 * time.Second
	maxRetries     = 3
	retryBaseDelay = 500 * time.Millisecond
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type Config struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

func NewClient(config Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}

	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	_, err := url.Parse(config.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Client{
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) do(req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ApiErrorResponse
		if json.Unmarshal(body, &errResp) == nil && (errResp.Message != "" || len(errResp.Errors) > 0) {
			if len(errResp.Errors) > 0 {
				return fmt.Errorf("API error (status %d): %v", resp.StatusCode, errResp.Errors)
			}
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Message)
		}
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) doWithRetry(req *http.Request, result interface{}) error {
	var bodyBytes []byte
	if req.GetBody != nil {
		getBody := req.GetBody
		if b, err := getBody(); err == nil {
			bodyBytes, _ = io.ReadAll(b)
			b.Close()
		}
	} else if req.Body != nil {
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bodyBytes)), nil
		}
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryBaseDelay * time.Duration(1<<(attempt-1))
			time.Sleep(delay)
		}

		retryReq := req.Clone(req.Context())
		if bodyBytes != nil && (retryReq.Method == http.MethodPost || retryReq.Method == http.MethodPut) {
			retryReq.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		lastErr = c.do(retryReq, result)
		if lastErr == nil {
			return nil
		}

		if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodDelete {
			if !isRetryableError(lastErr) {
				return lastErr
			}
		}
	}
	return lastErr
}

func isRetryableError(err error) bool {
	errMsg := err.Error()
	return containsAny(errMsg, "status 429", "status 500", "status 502", "status 503", "status 504", "connection refused", "timeout")
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}

func (c *Client) GetAccount(ctx context.Context) (*Account, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/account", nil)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := c.doWithRetry(req, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *Client) GetLimits(ctx context.Context) (*LimitsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/cloud/limits", nil)
	if err != nil {
		return nil, err
	}

	var limits LimitsResponse
	if err := c.doWithRetry(req, &limits); err != nil {
		return nil, err
	}

	return &limits, nil
}

func (c *Client) GetProjects(ctx context.Context) (*ProjectsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/project", nil)
	if err != nil {
		return nil, err
	}

	var projects ProjectsResponse
	if err := c.doWithRetry(req, &projects); err != nil {
		return nil, err
	}

	return &projects, nil
}

func (c *Client) GetOffers(ctx context.Context) (*OffersResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/cloud/offers", nil)
	if err != nil {
		return nil, err
	}

	var offers OffersResponse
	if err := c.doWithRetry(req, &offers); err != nil {
		return nil, err
	}

	return &offers, nil
}

func (c *Client) GetSystems(ctx context.Context) (*SystemsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/cloud/systems", nil)
	if err != nil {
		return nil, err
	}

	var systems SystemsResponse
	if err := c.doWithRetry(req, &systems); err != nil {
		return nil, err
	}

	return &systems, nil
}

func (c *Client) ListCloudServers(ctx context.Context) (*CloudServersResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/cloud-server", nil)
	if err != nil {
		return nil, err
	}

	var servers CloudServersResponse
	if err := c.doWithRetry(req, &servers); err != nil {
		return nil, err
	}

	return &servers, nil
}

func (c *Client) GetCloudServer(ctx context.Context, uuid string) (*CloudServer, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/cloud-server/%s", uuid), nil)
	if err != nil {
		return nil, err
	}

	var response CloudServerResponse
	if err := c.doWithRetry(req, &response); err != nil {
		return nil, err
	}

	return &response.CloudServer, nil
}

func (c *Client) CreateCloudServer(ctx context.Context, createReq CreateCloudServerRequest) (*CreateCloudServerResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/cloud/order", createReq)
	if err != nil {
		return nil, err
	}

	var result CreateCloudServerResponse
	if err := c.doWithRetry(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeleteCloudServer(ctx context.Context, uuid string) (*CloudServerActionResponse, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/cloud-server/%s", uuid), nil)
	if err != nil {
		return nil, err
	}

	var result CloudServerActionResponse
	if err := c.doWithRetry(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) StartCloudServer(ctx context.Context, uuid string) (*CloudServerActionResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/cloud-server/%s/start", uuid), nil)
	if err != nil {
		return nil, err
	}

	var result CloudServerActionResponse
	if err := c.doWithRetry(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) StopCloudServer(ctx context.Context, uuid string) (*CloudServerActionResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/cloud-server/%s/stop", uuid), nil)
	if err != nil {
		return nil, err
	}

	var result CloudServerActionResponse
	if err := c.doWithRetry(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) RestartCloudServer(ctx context.Context, uuid string) (*CloudServerActionResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/cloud-server/%s/restart", uuid), nil)
	if err != nil {
		return nil, err
	}

	var result CloudServerActionResponse
	if err := c.doWithRetry(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) ResizeCloudServer(ctx context.Context, uuid string, resizeReq ResizeCloudServerRequest) (*CloudServerActionResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/cloud-server/%s/resize", uuid), resizeReq)
	if err != nil {
		return nil, err
	}

	var result CloudServerActionResponse
	if err := c.doWithRetry(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetCloudServerPassword(ctx context.Context, uuid, passwordType string) (*PasswordResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/cloud-server/%s/password?password=%s", uuid, passwordType), nil)
	if err != nil {
		return nil, err
	}

	var result PasswordResponse
	if err := c.doWithRetry(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) WaitForCloudServerStatus(ctx context.Context, uuid, targetStatus string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for server %s to reach status %q", uuid, targetStatus)
		case <-ticker.C:
			server, err := c.GetCloudServer(ctx, uuid)
			if err != nil {
				return fmt.Errorf("failed to check server status: %w", err)
			}
			if server.Status == targetStatus {
				return nil
			}
			if server.Status == "error" {
				return fmt.Errorf("server %s entered error state", uuid)
			}
		}
	}
}

func (c *Client) WaitForCloudServerSetup(ctx context.Context, uuid string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for server %s setup to complete", uuid)
		case <-ticker.C:
			server, err := c.GetCloudServer(ctx, uuid)
			if err != nil {
				return fmt.Errorf("failed to check server setup status: %w", err)
			}
			if server.IsSetup == "1" {
				return nil
			}
			if server.IsError == "1" {
				return fmt.Errorf("server %s entered error state during setup", uuid)
			}
		}
	}
}
