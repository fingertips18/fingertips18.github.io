package v1

import (
	"fmt"
	"log"

	"github.com/blackmagiqq/ga4"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/client"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
)

type AnalyticsRepository interface {
	PageView(pageView domain.PageView) error
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
func (r *analyticsRepository) PageView(pageView domain.PageView) error {
	if err := pageView.Validate(); err != nil {
		return fmt.Errorf("failed to validate page view: %w", err)
	}

	log.Printf("Sending page view event on page %s location %s\n", pageView.PageTitle, pageView.PageLocation)

	err := r.analyticsAPI.SendEvent(
		ga4.Event{
			Name: "page_view",
			Params: map[string]any{
				"page_location": pageView.PageLocation,
				"page_title":    pageView.PageTitle,
			},
		},
		ga4.ClientID(utils.GenerateKey()),
	)

	if err != nil {
		return fmt.Errorf("failed to send event: %w", err)
	}

	log.Printf("Successful page view: pageTitle=%s, pageLocation=%s", pageView.PageTitle, pageView.PageLocation)

	return nil
}
