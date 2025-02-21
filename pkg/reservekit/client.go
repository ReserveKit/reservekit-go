package reservekit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultHost    = "https://api.reservekit.io"
	defaultVersion = "v1"
)

// Client represents a ReserveKit API client
type Client struct {
	secretKey  string
	host       string
	version    string
	service    *Service
	httpClient *http.Client
}

// NewClient creates a new ReserveKit client
func NewClient(secretKey string, opts ...Option) *Client {
	c := &Client{
		secretKey: secretKey,
		host:      defaultHost,
		version:   defaultVersion,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
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

	var bodyReader *bytes.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
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
