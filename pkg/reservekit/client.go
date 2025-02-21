package reservekit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultHost    = "https://api.reservekit.io"
	defaultVersion = "v1"
)

// Client represents the ReserveKit API client
type Client struct {
	secretKey  string
	host       string
	version    string
	httpClient *http.Client
	service    *Service
}

// NewClient creates a new ReserveKit API client
func NewClient(secretKey string, opts ...Option) *Client {
	c := &Client{
		secretKey:  secretKey,
		host:       defaultHost,    // default host
		version:    defaultVersion, // default version
		httpClient: &http.Client{}, // default http client
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Service returns the current service client
func (c *Client) Service() *Service {
	return c.service
}

// InitService initializes a service by ID
func (c *Client) InitService(serviceID int) error {
	service, err := c.getService(serviceID)
	if err != nil {
		return fmt.Errorf("failed to initialize service: %w", err)
	}

	c.service = service
	return nil
}

func (c *Client) request(method, path string, body interface{}, result interface{}) error {
	url := fmt.Sprintf("%s/%s%s", c.host, c.version, path)

	var reqBody io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return &apiErr
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) getService(serviceID int) (*Service, error) {
	path := fmt.Sprintf("/services/%d", serviceID)
	var result struct {
		Data ServiceData `json:"data"`
	}

	if err := c.request("GET", path, nil, &result); err != nil {
		return nil, err
	}

	return NewService(c, result.Data), nil
}
