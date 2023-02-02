package amplitude

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/pokt-foundation/utils-go/mock-client"
	"github.com/stretchr/testify/require"
)

func TestClient_AmplitudeClient(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("dummy", nil)

	for i := 0; i < 3; i++ {
		err := client.LogEvent(&Event{UserID: "1"})
		c.NoError(err)
	}

	mock.AddMultipleMockedPlainResponses(http.MethodPost, serverURL, []int{
		http.StatusRequestEntityTooLarge,
		http.StatusOK,
		http.StatusOK,
	}, []string{"ok", "ok", "ok"})

	err := client.Flush()
	c.NoError(err)
	c.Len(client.events, 0)
}

func TestClient_AmplitudeClientError(t *testing.T) {
	c := require.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient("dummy", &ClientOptions{
		Timeout: 2 * time.Minute,
	})

	err := client.Flush()
	c.Equal(err, ErrNoEvents)

	err = client.LogEvent(&Event{})
	c.Equal(err, ErrInvalidEvent)

	for i := 0; i < 3; i++ {
		err := client.LogEvent(&Event{UserID: "1"})
		c.NoError(err)
	}

	mock.AddMockedResponse(http.MethodPost, serverURL, http.StatusBadRequest, "not ok")

	err = client.Flush()
	c.EqualError(err, "error uploading events. Status Code: 400")

	mock.AddMultipleMockedPlainResponses(http.MethodPost, serverURL, []int{
		http.StatusRequestEntityTooLarge,
		http.StatusOK,
		http.StatusRequestEntityTooLarge,
	}, []string{"ok", "ok", "not ok"})

	err = client.Flush()
	c.Equal(err, ErrUnmanagedPayloadTooLarge)

	mock.AddMultipleMockedPlainResponses(http.MethodPost, serverURL, []int{
		http.StatusRequestEntityTooLarge,
		http.StatusRequestEntityTooLarge,
		http.StatusRequestEntityTooLarge,
		http.StatusRequestEntityTooLarge,
	}, []string{"not ok", "not ok", "not ok", "not ok"})

	err = client.Flush()
	c.Equal(err, ErrTooManyParts)
}
