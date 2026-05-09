package handlers

import (
	"log"
	"time"

	"github.com/SaranHiruthikM/newsletter-system/internal/domain"
	"github.com/SaranHiruthikM/newsletter-system/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NewsletterHandler struct {
	subRepo  repository.SubscriberRepository
	newsRepo repository.NewsletterRepository
}

type NewsletterSendRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewNewsletterHandler(subRepo repository.SubscriberRepository, newsRepo repository.NewsletterRepository) *NewsletterHandler {
	return &NewsletterHandler{subRepo: subRepo, newsRepo: newsRepo}
}

func (h *NewsletterHandler) Handle(c *fiber.Ctx) error {
	var req *NewsletterSendRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "error in parsing request"})
	}

	if req.Body == "" || req.Subject == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{"error": "either body or subject is missing"},
		)
	}

	newsletter := &domain.NewsletterSend{
		ID:        uuid.New().String(),
		Subject:   req.Subject,
		Body:      req.Body,
		Status:    domain.StatusPending,
		SentCount: 0,
		FailCount: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.newsRepo.Create(newsletter); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "something went wrong"},
		)
	}

	subscribers, err := h.subRepo.FindAllConfirmed()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			fiber.Map{"error": "something went wrong"},
		)
	}

	for _, subs := range subscribers {
		log.Printf("would send to: %s", subs.Email)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "newsletter dispatch started", "total": len(subscribers)})

}
