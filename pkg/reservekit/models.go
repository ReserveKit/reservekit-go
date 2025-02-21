package reservekit

import "time"

// TimeSlot represents a service time slot
type TimeSlot struct {
	ID          int       `json:"id"`
	ServiceID   int       `json:"service_id"`
	DayOfWeek   int       `json:"day_of_week"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	MaxBookings int       `json:"max_bookings"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BookingRequest represents a booking creation request
type BookingRequest struct {
	CustomerName  string    `json:"customer_name,omitempty"`
	CustomerEmail string    `json:"customer_email,omitempty"`
	CustomerPhone string    `json:"customer_phone,omitempty"`
	Date          time.Time `json:"date"`
	TimeSlotID    int       `json:"time_slot_id"`
}

// Booking represents a service booking
type Booking struct {
	ID         int       `json:"id"`
	ServiceID  int       `json:"service_id"`
	CustomerID int       `json:"customer_id"`
	TimeSlotID int       `json:"time_slot_id"`
	Date       time.Time `json:"date"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
