package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/client"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
)

type EmailRepository interface {
	Send(send domain.SendEmail) error
}

type EmailRepositoryConfig struct {
	ServiceID      string            `json:"service_id"`
	TemplateID     string            `json:"template_id"`
	UserID         string            `json:"user_id"`
	AccessToken    string            `json:"accessToken,omitempty"`
	TemplateParams map[string]string `json:"template_params"`

	httpAPI client.HttpAPI
}

type emailRepository struct {
	payload EmailRepositoryConfig
	httpAPI client.HttpAPI
}

// NewEmailRepository creates and returns a new instance of EmailRepository using the provided configuration.
// It initializes the repository with the specified EmailRepositoryConfig.
func NewEmailRepository(cfg EmailRepositoryConfig) EmailRepository {
	payload := EmailRepositoryConfig{
		ServiceID:      cfg.ServiceID,
		TemplateID:     cfg.TemplateID,
		UserID:         cfg.UserID,
		AccessToken:    cfg.AccessToken,
		TemplateParams: cfg.TemplateParams,
	}

	if payload.TemplateParams == nil {
		payload.TemplateParams = make(map[string]string)
	}

	httpAPI := cfg.httpAPI
	if httpAPI == nil {
		httpAPI = client.NewHTTPAPI()
	}

	return &emailRepository{
		payload: payload,
		httpAPI: httpAPI,
	}
}

// Send sends an email using the EmailJS API with the provided name, email, subject, and message.
// It updates the template parameters with the given values, marshals the payload to JSON,
// and performs an HTTP POST request to the EmailJS endpoint.
// Returns an error if the request fails or if the response status is not OK.
func (r *emailRepository) Send(send domain.SendEmail) error {
	if err := send.Validate(); err != nil {
		return fmt.Errorf("failed to validate send: %w", err)
	}

	log.Printf("Sending email via EmailJS")

	params := make(map[string]string, len(r.payload.TemplateParams)+4)
	maps.Copy(params, r.payload.TemplateParams)
	params["name"] = send.Name
	params["email"] = send.Email
	params["subject"] = send.Subject
	params["message"] = send.Message

	payload := r.payload
	payload.TemplateParams = params

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.emailjs.com/api/v1.0/email/send", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.httpAPI.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send: [status=%s,message=%s]", resp.Status, string(respBody))
	}

	log.Printf("Email sent successfully via EmailJS")

	return nil
}
