package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	mockRepo "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1/mocks"
	"github.com/stretchr/testify/assert"
)

type analyticsHandlerTestFixture struct {
	t                 *testing.T
	mockAnalyticsRepo *mockRepo.MockAnalyticsRepository
	analyticsHandler  AnalyticsHandler
}

func newAnalyticsHandlerTestFixture(t *testing.T) *analyticsHandlerTestFixture {
	mockAnalyticsRepo := new(mockRepo.MockAnalyticsRepository)

	analyticsHandler := NewAnalyticsServiceHandler(
		AnalyticsServiceConfig{
			analyticsRepo: mockAnalyticsRepo,
		},
	)

	return &analyticsHandlerTestFixture{
		t:                 t,
		mockAnalyticsRepo: mockAnalyticsRepo,
		analyticsHandler:  analyticsHandler,
	}
}

func TestAnalyticsServiceHandler_PageView(t *testing.T) {
	validReq := domain.PageView{
		PageLocation: "http://example.com/home",
		PageTitle:    "Homepage",
	}
	validBody, _ := json.Marshal(validReq)

	type Given struct {
		method   string
		body     string
		mockRepo func(m *mockRepo.MockAnalyticsRepository)
	}
	type Expected struct {
		code int
		body string
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"success": {
			given: Given{
				method: http.MethodPost,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockAnalyticsRepository) {
					m.EXPECT().
						PageView(validReq).
						Return(nil)
				},
			},
			expected: Expected{
				code: http.StatusOK,
				body: toJSON(map[string]string{
					"message":      "Page view recorded successfully",
					"pageLocation": validReq.PageLocation,
					"pageTitle":    validReq.PageTitle,
				}),
			},
		},
		"invalid method": {
			given: Given{
				method: http.MethodGet,
				body:   "",
			},
			expected: Expected{
				code: http.StatusMethodNotAllowed,
				body: "Method not allowed: only POST is supported\n",
			},
		},
		"invalid json": {
			given: Given{
				method: http.MethodPost,
				body:   `{"pageLocation":}`,
			},
			expected: Expected{
				code: http.StatusBadRequest,
				body: "Invalid JSON in request body\n",
			},
		},
		"repo error": {
			given: Given{
				method: http.MethodPost,
				body:   string(validBody),
				mockRepo: func(m *mockRepo.MockAnalyticsRepository) {
					m.EXPECT().
						PageView(validReq).
						Return(errors.New("tracking failed"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to view page: tracking failed\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newAnalyticsHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockAnalyticsRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/analytics/page-view", strings.NewReader(tt.given.body))
			w := httptest.NewRecorder()

			f.analyticsHandler.(*analyticsServiceHandler).PageView(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "{") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockAnalyticsRepo.AssertExpectations(t)
		})
	}
}

func TestAnalyticsServiceHandler_PageView_Routing(t *testing.T) {
	validReq := domain.PageView{
		PageLocation: "http://example.com/home",
		PageTitle:    "Homepage",
	}
	validBody, _ := json.Marshal(validReq)

	expectedResp, _ := json.Marshal(map[string]string{
		"message":      "Page view recorded successfully",
		"pageLocation": validReq.PageLocation,
		"pageTitle":    validReq.PageTitle,
	})

	f := newAnalyticsHandlerTestFixture(t)

	// Mock expectation
	f.mockAnalyticsRepo.EXPECT().
		PageView(validReq).
		Return(nil)

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/analytics/page-view", bytes.NewReader(validBody))
	w := httptest.NewRecorder()

	// Verify handler implements http.Handler
	handler, ok := f.analyticsHandler.(http.Handler)
	assert.True(t, ok, "analyticsHandler should implement http.Handler")

	// Route through ServeHTTP
	handler.ServeHTTP(w, req)

	// Validate response
	res := w.Result()
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, string(expectedResp), string(body))

	f.mockAnalyticsRepo.AssertExpectations(t)
}
