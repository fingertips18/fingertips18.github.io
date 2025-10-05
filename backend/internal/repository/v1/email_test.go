package v1

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	client "github.com/fingertips18/fingertips18.github.io/backend/internal/client/mocks"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type emailRepositoryTestFixture struct {
	t               *testing.T
	mockHttpAPI     *client.MockHttpAPI
	emailRepository emailRepository
}

func newEmailRepositoryTestFixture(t *testing.T) *emailRepositoryTestFixture {
	mockHttpAPI := new(client.MockHttpAPI)
	emailRepository := &emailRepository{
		payload: EmailRepositoryConfig{
			ServiceID:      "service_xxx",
			TemplateID:     "template_xxx",
			UserID:         "user_xxx",
			TemplateParams: make(map[string]string),
		},
		httpAPI: mockHttpAPI,
	}

	return &emailRepositoryTestFixture{
		t:               t,
		mockHttpAPI:     mockHttpAPI,
		emailRepository: *emailRepository,
	}
}

func TestEmailRepository_Send(t *testing.T) {
	httpErr := errors.New("http error")
	serverErr := errors.New("server error")
	missingNameErr := errors.New("failed to validate send: name missing")
	missingEmailErr := errors.New("failed to validate send: email missing")
	missingSubjectErr := errors.New("failed to validate send: subject missing")

	payload := domain.SendEmail{
		Name:    "Test User",
		Email:   "test@user.com",
		Subject: "Test subject",
		Message: "Test message...",
	}

	type Given struct {
		payload  domain.SendEmail
		mockSend func(m *client.MockHttpAPI)
	}

	type Expected struct {
		err error
	}

	tests := map[string]struct {
		given    Given
		expected Expected
	}{
		"Successful send": {
			given: Given{
				payload: payload,
				mockSend: func(m *client.MockHttpAPI) {
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(`ok`)),
						}, nil)
				},
			},
			expected: Expected{
				err: nil,
			},
		},
		"HTTP client error": {
			given: Given{
				payload: payload,
				mockSend: func(m *client.MockHttpAPI) {
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 500,
						}, httpErr)
				},
			},
			expected: Expected{
				err: fmt.Errorf("failed to send HTTP request: %w", httpErr),
			},
		},
		"Non-200 response": {
			given: Given{
				payload: payload,
				mockSend: func(m *client.MockHttpAPI) {
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 500,
							Status:     "500 Internal Server Error",
							Body:       io.NopCloser(bytes.NewBufferString(serverErr.Error())),
						}, nil)
				},
			},
			expected: Expected{
				err: fmt.Errorf("failed to send: [status=500 Internal Server Error,message=%s]", serverErr),
			},
		},
		"Missing name": {
			given: Given{
				payload: domain.SendEmail{
					Name:    "",
					Email:   payload.Email,
					Subject: payload.Subject,
					Message: payload.Message,
				},
				mockSend: nil,
			},
			expected: Expected{
				err: missingNameErr,
			},
		},
		"Missing email": {
			given: Given{
				payload: domain.SendEmail{
					Name:    payload.Name,
					Email:   "",
					Subject: payload.Subject,
					Message: payload.Message,
				},
				mockSend: nil,
			},
			expected: Expected{
				err: missingEmailErr,
			},
		},
		"Missing subject": {
			given: Given{
				payload: domain.SendEmail{
					Name:    payload.Name,
					Email:   payload.Email,
					Subject: "",
					Message: payload.Message,
				},
				mockSend: nil,
			},
			expected: Expected{
				err: missingSubjectErr,
			},
		},
		"Missing message": {
			given: Given{
				payload: domain.SendEmail{
					Name:    payload.Name,
					Email:   payload.Email,
					Subject: payload.Subject,
					Message: "",
				},
				mockSend: func(m *client.MockHttpAPI) {
					m.EXPECT().Do(mock.AnythingOfType("*http.Request")).
						Return(&http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(`ok`)),
						}, nil)
				},
			},
			expected: Expected{
				err: nil,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			f := newEmailRepositoryTestFixture(t)

			if test.given.mockSend != nil {
				test.given.mockSend(f.mockHttpAPI)
			}

			err := f.emailRepository.Send(test.given.payload)

			if test.expected.err != nil {
				assert.EqualError(t, err, test.expected.err.Error())
			} else {
				assert.NoError(t, err)
			}

			f.mockHttpAPI.AssertExpectations(t)
		})
	}
}
