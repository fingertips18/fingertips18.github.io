package v1

import (
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

type emailHandlerTestFixture struct {
	t             *testing.T
	mockEmailRepo *mockRepo.MockEmailRepository
	emailHandler  EmailHandler
}

func newEmailHandlerTestFixture(t *testing.T) *emailHandlerTestFixture {
	mockEmailRepo := new(mockRepo.MockEmailRepository)

	emailHandler := NewEmailServiceHandler(
		EmailServiceConfig{
			emailRepo: mockEmailRepo,
		},
	)

	return &emailHandlerTestFixture{
		t:             t,
		mockEmailRepo: mockEmailRepo,
		emailHandler:  emailHandler,
	}
}

func TestEmailServiceHandler_Send(t *testing.T) {
	validReq := domain.SendEmail{
		Name:    "John Doe",
		Email:   "john@example.com",
		Message: "Hello, world!",
	}
	validBody, _ := json.Marshal(validReq)

	type Given struct {
		method   string
		body     string
		mockRepo func(m *mockRepo.MockEmailRepository)
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
				mockRepo: func(m *mockRepo.MockEmailRepository) {
					m.EXPECT().
						Send(validReq).
						Return(nil)
				},
			},
			expected: Expected{
				code: http.StatusAccepted,
				body: toJSON(map[string]string{
					"message": "Email sent successfully",
					"email":   validReq.Email,
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
				body:   `{"email": "invalid",}`,
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
				mockRepo: func(m *mockRepo.MockEmailRepository) {
					m.EXPECT().
						Send(validReq).
						Return(errors.New("smtp error"))
				},
			},
			expected: Expected{
				code: http.StatusInternalServerError,
				body: "Failed to send email: smtp error\n",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			f := newEmailHandlerTestFixture(t)

			if tt.given.mockRepo != nil {
				tt.given.mockRepo(f.mockEmailRepo)
			}

			req := httptest.NewRequest(tt.given.method, "/email/send", strings.NewReader(tt.given.body))
			w := httptest.NewRecorder()

			f.emailHandler.(*emailServiceHandler).Send(w, req)

			res := w.Result()
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, tt.expected.code, res.StatusCode)

			if strings.HasPrefix(tt.expected.body, "{") {
				assert.JSONEq(t, tt.expected.body, string(body))
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			f.mockEmailRepo.AssertExpectations(t)
		})
	}
}
