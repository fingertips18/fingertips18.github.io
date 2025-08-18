package client

import "github.com/blackmagiqq/ga4"

type GoogleAnalyticsAPI interface {
	SendEvent(event ga4.Event, clientID ga4.ClientID) error
}

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
