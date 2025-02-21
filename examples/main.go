package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ReserveKit/reservekit-go/pkg/reservekit"
)

func main() {
	// Create a new client
	client := reservekit.NewClient("your-api-key")

	// Initialize a service
	err := client.InitService(1)
	if err != nil {
		log.Fatal(err)
	}

	// Get available time slots
	slots, err := client.Service().GetTimeSlots()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d time slots\n", len(slots))

	// Create a booking
	booking, err := client.Service().CreateBooking(&reservekit.BookingRequest{
		CustomerName:  "John Doe",
		CustomerEmail: "john@example.com",
		CustomerPhone: "+1234567890",
		Date:          time.Now().AddDate(0, 0, 1),
		TimeSlotID:    slots[0].ID,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created booking with ID: %d\n", booking.ID)
}
