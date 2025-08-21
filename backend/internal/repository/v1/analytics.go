package v1

import (
	"fmt"
	"log"

	"github.com/blackmagiqq/ga4"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/client"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
)

type AnalyticsRepository interface {
	PageView(pageLocation, pageTitle string) error
}

type AnalyticsRepositoryConfig struct {
	GoogleMeasurementID string
	GoogleAPISecret     string

	analyticsAPI client.GoogleAnalyticsAPI
}

type analyticsRepository struct {
	GoogleMeasurementID string
	GoogleAPISecret     string

	analyticsAPI client.GoogleAnalyticsAPI
}

func NewAnalyticsRepository(cfg AnalyticsRepositoryConfig) AnalyticsRepository {
	analyticsAPI := cfg.analyticsAPI

	if analyticsAPI == nil {
		analyticsAPI = client.NewGoogleAnalyticsAPI(
			cfg.GoogleMeasurementID,
			cfg.GoogleAPISecret,
		)
	}

	return &analyticsRepository{
		analyticsAPI: analyticsAPI,
	}
}

// PageView sends a "page_view" event to the analytics API with the specified page location and title.
// It logs the page visit and returns an error if the event could not be sent.
//
// Parameters:
//   - pageLocation: The URL or identifier of the page being viewed.
//   - pageTitle: The title of the page being viewed.
//
// Returns:
//   - error: An error if sending the analytics event fails, otherwise nil.
func (r *analyticsRepository) PageView(pageLocation, pageTitle string) error {
	log.Printf("Page %s visited", pageLocation)

	err := r.analyticsAPI.SendEvent(
		ga4.Event{
			Name: "page_view",
			Params: map[string]interface{}{
				"page_location": pageLocation,
				"page_title":    pageTitle,
			},
		},
		ga4.ClientID(utils.GenerateKey()),
	)

	if err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

	return nil
}
