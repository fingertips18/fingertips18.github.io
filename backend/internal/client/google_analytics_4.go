package client

import "github.com/blackmagiqq/ga4"

type GoogleAnalyticsAPI interface {
	SendEvent(event ga4.Event, clientID ga4.ClientID) error
}

// NewGoogleAnalyticsAPI creates and returns a new instance of GoogleAnalyticsAPI using the provided
// measurementID and apiSecret. It initializes the GA4 client and panics if there is an error during
// client creation.
//
// Parameters:
//   - measurementID: The Google Analytics 4 Measurement ID.
//   - apiSecret: The API secret associated with the Measurement ID.
//
// Returns:
//   - GoogleAnalyticsAPI: An initialized Google Analytics API client.
func NewGoogleAnalyticsAPI(measurementID, apiSecret string) GoogleAnalyticsAPI {
	client, err := ga4.NewGA4Client(
		measurementID,
		apiSecret,
		true,
	)

	if err != nil {
		panic(err)
	}

	return client
}
