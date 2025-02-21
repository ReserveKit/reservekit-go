package integration

import (
	"testing"

	"github.com/ReserveKit/reservekit-go/pkg/reservekit"
	"github.com/stretchr/testify/assert"
)

func TestClientIntegration(t *testing.T) {
	client := reservekit.NewClient(
		"test-key",
		reservekit.WithHost("https://api.staging.reservekit.io"),
	)

	err := client.InitService(1)
	assert.NoError(t, err)
	assert.NotNil(t, client.Service())
	assert.Equal(t, 1, client.Service().ID)
}
