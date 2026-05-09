package handlers

import (
	"log"
	"time"

	"github.com/SaranHiruthikM/newsletter-system/internal/domain"
	"github.com/SaranHiruthikM/newsletter-system/internal/queue"
	"github.com/SaranHiruthikM/newsletter-system/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SubscribeHandler struct {
	repo repository.SubscriberRepository
	pub  *queue.Publisher
}

type SubscribeRequest struct {
	Email string `json:"email"`
}

func NewSubscriberHandler(repo repository.SubscriberRepository, publisher *queue.Publisher) *SubscribeHandler {
	return &SubscribeHandler{
		repo: repo,
		pub:  publisher,
	}
}

func (h *SubscribeHandler) Handle(c *fiber.Ctx) error {
	var req SubscribeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to save subscriber"})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email id is required"})
	}

	sub, err := h.repo.FindByEmail(req.Email)

	if sub != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already subscribed"})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "something went wrong"})
	}

	token := uuid.New().String()

	subscriber := &domain.Subscriber{
		ID:             uuid.New().String(),
		Confirmed:      false,
		Token:          token,
		Email:          req.Email,
		TokenExpiresAt: time.Now().Add(time.Hour * 24),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = h.repo.Create(subscriber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create a subscriber"})
	}

	log.Printf("confirmation token for %s: %s", req.Email, token)

	h.pub.PublishConfirmation(queue.ConfirmationPayload{
		Email: subscriber.Email,
		Token: subscriber.Token,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "subscription successful, please check your mail"})
}
