package queue

type ConfirmationPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type NewsletterPayload struct {
	NewsletterID string `json:"newsletter_id"`
	Email        string `json:"email"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
}
