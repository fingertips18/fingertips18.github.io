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

// PageView validates the provided domain.PageView and sends a "page_view" event to the configured analytics API.
// The event includes "page_location" and "page_title" parameters and is sent using a generated client ID.
// It logs the attempt and the successful send. Returns an error if validation or sending the event fails.
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
