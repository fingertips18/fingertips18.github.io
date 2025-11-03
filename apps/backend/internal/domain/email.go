package domain

import "errors"

type SendEmail struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (s SendEmail) Validate() error {
	if s.Name == "" {
		return errors.New("name missing")
	}
	if s.Email == "" {
		return errors.New("email missing")
	}
	if s.Subject == "" {
		return errors.New("subject missing")
	}

	return nil
}
