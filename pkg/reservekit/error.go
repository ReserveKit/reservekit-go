package reservekit

import "fmt"

// APIError represents an error returned by the API
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %s (code: %s, status: %d)", e.Message, e.Code, e.Status)
}
