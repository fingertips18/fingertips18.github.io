package v1

import (
	"errors"
	"fmt"
	"testing"

	client "github.com/fingertips18/fingertips18.github.io/backend/internal/client/mocks"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type analyticsRepositoryTestFixture struct {
	t                   *testing.T
	mockAnalyticsAPI    *client.MockGoogleAnalyticsAPI
	analyticsRepository analyticsRepository
}

func newAnalyticsRepositoryTestFixture(t *testing.T) *analyticsRepositoryTestFixture {
	mockAnalyticsAPI := new(client.MockGoogleAnalyticsAPI)
	analyticsRepo := &analyticsRepository{
		analyticsAPI: mockAnalyticsAPI,
	}

	return &analyticsRepositoryTestFixture{
		t:                   t,
		mockAnalyticsAPI:    mockAnalyticsAPI,
		analyticsRepository: *analyticsRepo,
	}
}

func TestAnalyticsRepository_PageView(t *testing.T) {
	eventErr := errors.New("Event error")
	missingPageLocationErr := errors.New("failed to validate page view: pageLocation missing")
	missingPageTitleErr := errors.New("failed to validate page view: pageTitle missing")

	type Given struct {
		pageView     domain.PageView
		mockPageView func(m *client.MockGoogleAnalyticsAPI)
	}

	type Expected struct {
		err error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful page view": {
			given: Given{
				pageView: domain.PageView{
					PageLocation: "/page",
					PageTitle:    "Page",
				},
				mockPageView: func(m *client.MockGoogleAnalyticsAPI) {
					m.EXPECT().SendEvent(
						mock.AnythingOfType("ga4.Event"),
						mock.AnythingOfType("ga4.ClientID"),
					).Return(nil)
				},
			},
			expected: Expected{
				err: nil,
			},
		},
		"SendEvent returns error": {
			given: Given{
				pageView: domain.PageView{
					PageLocation: "/error-Page",
					PageTitle:    "Error Page",
				},
				mockPageView: func(m *client.MockGoogleAnalyticsAPI) {
					m.EXPECT().SendEvent(
						mock.AnythingOfType("ga4.Event"),
						mock.AnythingOfType("ga4.ClientID"),
					).Return(eventErr)
				},
			},
			expected: Expected{
				err: fmt.Errorf("failed to send event: %w", eventErr),
			},
		},
		"Missing page location": {
			given: Given{
				pageView: domain.PageView{
					PageLocation: "",
					PageTitle:    "Page",
				},
				mockPageView: nil,
			},
			expected: Expected{
				err: missingPageLocationErr,
			},
		},
		"Missing page title": {
			given: Given{
				pageView: domain.PageView{
					PageLocation: "/page",
					PageTitle:    "",
				},
				mockPageView: nil,
			},
			expected: Expected{
				err: missingPageTitleErr,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newAnalyticsRepositoryTestFixture(t)

			if test.given.mockPageView != nil {
				test.given.mockPageView(f.mockAnalyticsAPI)
			}

			err := f.analyticsRepository.PageView(test.given.pageView)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Empty(t, err)
			}

			f.mockAnalyticsAPI.AssertExpectations(t)
		})
	}
}
