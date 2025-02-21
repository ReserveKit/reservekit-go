# ReserveKit Go SDK

A Go client library for the ReserveKit API.

## Installation

```bash
go get github.com/ReserveKit/reservekit-go
```

## Quick Start

```go
package main
import (
"fmt"
"github.com/ReserveKit/reservekit-go/pkg/reservekit"
)
func main() {
// Create a new client
client := reservekit.NewClient("your-api-key")
// Initialize a service
err := client.InitService(1)
if err != nil {
panic(err)
}
// Get available time slots
slots, err := client.Service().GetTimeSlots()
if err != nil {
panic(err)
}
fmt.Printf("Found %d time slots\n", len(slots))
}
```

## Documentation

For more detailed information on the API and available methods, please refer to
the [ReserveKit API Documentation](https://docs.reservekit.io).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file
for details.
