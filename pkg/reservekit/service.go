package reservekit

import (
	"fmt"
	"time"
)

// Service represents a ReserveKit service client
type Service struct {
	client      *Client
	ID          int       `json:"id"`
	ProviderID  string    `json:"provider_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Timezone    string    `json:"timezone"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ServiceData represents the raw service data from the API
type ServiceData struct {
	ID          int       `json:"id"`
	ProviderID  string    `json:"provider_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Timezone    string    `json:"timezone"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewService creates a new service client
func NewService(client *Client, data ServiceData) *Service {
	return &Service{
		client:      client,
		ID:          data.ID,
		ProviderID:  data.ProviderID,
		Name:        data.Name,
		Description: data.Description,
		Timezone:    data.Timezone,
		Version:     data.Version,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}

// GetTimeSlots retrieves available time slots for the service
func (s *Service) GetTimeSlots() ([]TimeSlot, error) {
	var result struct {
		Data struct {
			TimeSlots []TimeSlot `json:"time_slots"`
		} `json:"data"`
	}

	err := s.client.request(
		"GET",
		fmt.Sprintf("/services/%d/time-slots", s.ID),
		nil,
		&result,
	)
	if err != nil {
		return nil, err
	}

	return result.Data.TimeSlots, nil
}

// CreateBooking creates a new booking for the service
func (s *Service) CreateBooking(req *BookingRequest) (*Booking, error) {
	var result struct {
		Data Booking `json:"data"`
	}

	err := s.client.request(
		"POST",
		fmt.Sprintf("/services/%d/bookings", s.ID),
		req,
		&result,
	)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}
