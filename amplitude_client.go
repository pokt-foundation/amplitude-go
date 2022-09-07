// Package amplitude is a collection of some functions to call Amplitude API
// Also manages the events sending by parts
package amplitude

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pokt-foundation/utils-go/client"
)

const (
	defaultRetries = 0
	defaultTimeout = 1 * time.Minute
)

var (
	serverURL = "https://api2.amplitude.com/2/httpapi"

	// ErrInvalidEvent error when event is not valid
	ErrInvalidEvent = errors.New("invalid event")
	// ErrUploadingEvents error when upload events API call fails
	ErrUploadingEvents = errors.New("error uploading events")
	// ErrTooManyParts error when events need to be divided by more than their size
	ErrTooManyParts = errors.New("error too many parts")
	// ErrUnmanagedPayloadTooLarge error when API throw too large in not the first request
	ErrUnmanagedPayloadTooLarge = errors.New("unmmanaged payload too large")

	// Errors managed by the SDK
	errPayloadTooLarge = errors.New("payload too large")
)

// Client struct handler for amplitude client
type Client struct {
	apiKey     string
	httpClient *client.Client
	events     []*Event
}

// ClientOptions struct handler for optional parameters in client's creation
type ClientOptions struct {
	Retries int
	Timeout time.Duration
}

// NewClient returns a new Client instance with given input
func NewClient(apiKey string, options *ClientOptions) *Client {
	if options == nil {
		return &Client{
			apiKey:     apiKey,
			httpClient: client.NewCustomClient(defaultRetries, defaultTimeout),
		}
	}

	return &Client{
		apiKey:     apiKey,
		httpClient: client.NewCustomClient(options.Retries, options.Timeout),
	}
}

// LogEvent adds event to client after validating it
func (c *Client) LogEvent(event *Event) error {
	if !event.IsValid() {
		return ErrInvalidEvent
	}

	c.events = append(c.events, event)

	return nil
}

// Flush uploads all logged events to Amplitude API
func (c *Client) Flush() error {
	parts := 1

	for {
		err := c.uploadEventsByParts(parts)
		if err != nil {
			if errors.Is(err, errPayloadTooLarge) {
				parts++

				continue
			}

			return err
		}

		break
	}

	c.events = nil

	return nil
}

func (c *Client) uploadEventsByParts(parts int) error {
	if parts > len(c.events) {
		return ErrTooManyParts
	}

	partSize := len(c.events) / parts
	modulus := len(c.events) % parts
	lastIndex := 0

	for i := 0; i < parts; i++ {
		fromIndex := lastIndex

		// adds the parts size as the last index to be uploaded
		toIndex := lastIndex + partSize

		// the first part sent should be the biggest one
		// to ensure it will return payload too large if it needs to
		if i < modulus {
			toIndex++
		}

		err := c.uploadEvents(c.events[fromIndex:toIndex])
		if err != nil {
			if errors.Is(err, errPayloadTooLarge) && i != 0 {
				// for now we will assume if the first parts is upload correctly everything should
				// in case it doesn't we'll return this error
				return ErrUnmanagedPayloadTooLarge
			}

			return err
		}

		lastIndex = toIndex
	}

	return nil
}

func (c *Client) uploadEvents(events []*Event) error {
	response, err := c.httpClient.PostWithURLJSONParams(serverURL, &uploadEventInput{
		APIKey: c.apiKey,
		Events: events,
	}, http.Header{})
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusRequestEntityTooLarge {
		return errPayloadTooLarge
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s. Status Code: %v", ErrUploadingEvents, response.StatusCode)
	}

	return nil
}
