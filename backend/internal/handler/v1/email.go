package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	v1 "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1"
)

type EmailService interface {
	http.Handler
	Send(w http.ResponseWriter, r *http.Request)
}

type EmailServiceConfig struct {
	ServiceID      string
	TemplateID     string
	UserID         string
	AccessToken    string
	TemplateParams map[string]string

	emailRepo v1.EmailRepository
}

type emailServiceHandler struct {
	emailRepo v1.EmailRepository
}

// NewEmailServiceHandler creates and returns a new instance of EmailService.
// It initializes the email repository using the provided EmailServiceConfig.
// If an email repository is not provided in the config, it creates a new one
// using the configuration parameters. The returned EmailService can be used
// to interact with email-related functionality.
func NewEmailServiceHandler(cfg EmailServiceConfig) EmailService {
	emailRepo := cfg.emailRepo
	if emailRepo == nil {
		emailRepo = v1.NewEmailRepository(
			v1.EmailRepositoryConfig{
				ServiceID:      cfg.ServiceID,
				TemplateID:     cfg.TemplateID,
				UserID:         cfg.UserID,
				AccessToken:    cfg.AccessToken,
				TemplateParams: cfg.TemplateParams,
			},
		)
	}

	return &emailServiceHandler{
		emailRepo: emailRepo,
	}
}

// ServeHTTP handles HTTP requests routed to the email service handler.
// It inspects the request path after trimming the "/email" prefix and dispatches
// the request to the appropriate handler method. If the path is "/send", it calls
// the Send method to process the request; otherwise, it responds with a 404 Not Found.
func (h *emailServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/email")

	switch path {
	case "/send":
		h.Send(w, r)
	default:
		http.NotFound(w, r)
	}
}

// Send handles HTTP POST requests to send an email.
// It expects a JSON payload in the request body containing the sender's name, email, and message.
// If the request method is not POST, it responds with "Method not allowed".
// On successful email sending, it responds with a JSON object {"status": "ok"} and HTTP 200 status.
// If there is an error decoding the request or sending the email, it responds with an appropriate HTTP error.
//
// @Security ApiKeyAuth
// @Summary Send an email
// @Description Sends an email with the provided details and returns a confirmation message.
// @Tags email
// @Accept json
// @Produce json
// @Param email body domain.SendEmail true "Email payload"
// @Success 202 {object} map[string]string "Confirmation message"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /email/send [post]
func (h *emailServiceHandler) Send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req domain.SendEmail
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	if err := h.emailRepo.Send(req.Name, req.Email, req.Subject, req.Message); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]string{
		"message": "Email sent successfully",
		"email":   req.Email,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(buf.Bytes())
}
