package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	repo "github.com/Fingertips18/fingertips18.github.io/backend/internal/repositories/email"
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

	emailRepo repo.EmailRepository
}

type emailServiceHandler struct {
	emailRepo repo.EmailRepository
}

type SendRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// NewEmailServiceHandler creates and returns a new instance of EmailService.
// It initializes the email repository using the provided EmailServiceConfig.
// If an email repository is not provided in the config, it creates a new one
// using the configuration parameters. The returned EmailService can be used
// to interact with email-related functionality.
func NewEmailServiceHandler(cfg EmailServiceConfig) EmailService {
	emailRepo := cfg.emailRepo
	if emailRepo == nil {
		emailRepo = repo.NewEmailRepository(repo.EmailRepositoryConfig{
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
func (h *emailServiceHandler) Send(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req SendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	if err := h.emailRepo.Send(req.Name, req.Email, req.Subject, req.Message); err != nil {
		http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]string{"status": "ok"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return
	}
}
