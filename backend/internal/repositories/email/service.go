package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type EmailRepository interface {
	Send(name, email, subject, message string) error
}

type EmailRepositoryConfig struct {
	ServiceID      string            `json:"service_id"`
	TemplateID     string            `json:"template_id"`
	UserID         string            `json:"user_id"`
	AccessToken    string            `json:"accessToken,omitempty"`
	TemplateParams map[string]string `json:"template_params"`
}

type emailRepository struct {
	payload EmailRepositoryConfig
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

	return &emailRepository{
		payload: payload,
	}
}

// Send sends an email using the EmailJS API with the provided name, email, subject, and message.
// It updates the template parameters with the given values, marshals the payload to JSON,
// and performs an HTTP POST request to the EmailJS endpoint.
// Returns an error if the request fails or if the response status is not OK.
func (r *emailRepository) Send(name, email, subject, message string) error {
	// Update template params with dynamic values before marshaling
	r.payload.TemplateParams["name"] = name
	r.payload.TemplateParams["email"] = email
	r.payload.TemplateParams["subject"] = subject
	r.payload.TemplateParams["message"] = message

	body, err := json.Marshal(r.payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	fmt.Println("Sending to EmailJS:", string(body))

	req, err := http.NewRequest("POST", "https://api.emailjs.com/api/v1.0/email/send", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("EmailJS error: %s - %s", resp.Status, string(respBody))
	}

	return nil
}
