package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ReserveKit/reservekit-go/pkg/reservekit"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/v1/services/1":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": {"id": 1, "name": "Test Service"}}`))
		case "/v1/time-slots?service_id=1":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": { "time_slots": [{"id": 1, "service_id": 1, "day_of_week": 1, "start_time": "2024-01-01T09:00:00Z", "end_time": "2024-01-01T10:00:00Z", "max_bookings": 1}]}}`))
		case "/v1/bookings?service_id=1":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": {"id": 1, "service_id": 1, "time_slot_id": 1, "status": "confirmed"}}`))
		case "/v1/services/99999":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"status": 404, "message": "Service not found", "code": "not_found"}`))
		default:
			fmt.Println("404 Request URL:", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"status": 404, "message": "Not found", "code": "not_found"}`))
		}
	}))
}

func TestClientIntegration(t *testing.T) {
	// Setup test server
	server := setupTestServer()
	defer server.Close()

	// Create client with test server URL
	client := reservekit.NewClient(
		"test-key",
		reservekit.WithHost(server.URL),
		reservekit.WithVersion("v1"), // Make sure version is set
	)

	t.Run("InitService", func(t *testing.T) {
		err := client.InitService(1)
		assert.NoError(t, err)
		assert.NotNil(t, client.Service())
		assert.Equal(t, 1, client.Service().ID)
	})

	t.Run("GetTimeSlots", func(t *testing.T) {
		// First initialize service
		err := client.InitService(1)
		assert.NoError(t, err)

		slots, err := client.Service().GetTimeSlots()
		assert.NoError(t, err)
		assert.NotEmpty(t, slots)

		// Verify time slot structure
		for _, slot := range slots {
			assert.NotZero(t, slot.ID)
			assert.Equal(t, 1, slot.ServiceID)
			assert.NotZero(t, slot.DayOfWeek)
			assert.False(t, slot.StartTime.IsZero())
			assert.False(t, slot.EndTime.IsZero())
			assert.True(t, slot.StartTime.Before(slot.EndTime))
		}
	})

	t.Run("CreateBooking", func(t *testing.T) {
		// First initialize service
		err := client.InitService(1)
		assert.NoError(t, err)

		// First get available slots
		slots, err := client.Service().GetTimeSlots()
		assert.NoError(t, err)
		assert.NotEmpty(t, slots)

		// Create a booking request
		bookingReq := &reservekit.BookingRequest{
			CustomerName:  "Test Customer",
			CustomerEmail: "test@example.com",
			CustomerPhone: "+1234567890",
			Date:          time.Now().AddDate(0, 0, 1), // Tomorrow
			TimeSlotID:    slots[0].ID,
		}

		booking, err := client.Service().CreateBooking(bookingReq)
		assert.NoError(t, err)
		assert.NotNil(t, booking)
		assert.NotZero(t, booking.ID)
		assert.Equal(t, 1, booking.ServiceID)
		assert.Equal(t, slots[0].ID, booking.TimeSlotID)
		assert.Equal(t, "confirmed", booking.Status)
	})

	t.Run("InitService_InvalidID", func(t *testing.T) {
		err := client.InitService(99999) // Using an invalid service ID
		assert.Error(t, err)
		var apiErr *reservekit.APIError
		assert.ErrorAs(t, err, &apiErr)
		assert.Equal(t, 404, apiErr.Status)
	})
}
