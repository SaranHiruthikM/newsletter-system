package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/SaranHiruthikM/newsletter-system/internal/config"
)

type ResendProvider struct {
	apiKey    string
	fromEmail string
	fromName  string
	baseURL   string
	client    *http.Client
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func NewResendProvider(cfg config.EmailConfig) *ResendProvider {
	log.Println("RESEND_API_KEY", cfg.ResendAPIKey)
	return &ResendProvider{
		apiKey:    cfg.ResendAPIKey,
		fromEmail: cfg.FromEmail,
		fromName:  cfg.FromName,
		baseURL:   cfg.ResendBaseURL,
		client:    &http.Client{Timeout: cfg.ResendTimeout},
	}

}

func (r *ResendProvider) Send(to, subject, body string) error {
	req := &resendRequest{
		From:    fmt.Sprintf("%s <%s>", r.fromName, r.fromEmail),
		To:      []string{to},
		Subject: subject,
		HTML:    body,
	}

	newBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	newReq, err := http.NewRequest("POST", r.baseURL+"/emails", bytes.NewBuffer(newBody))
	if err != nil {
		return err
	}

	newReq.Header.Set("Authorization", "Bearer "+r.apiKey)
	newReq.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(newReq)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("resend API error: status %d", resp.StatusCode)
	}

	return nil
}
