package handlers

import (
	"time"

	"github.com/SaranHiruthikM/newsletter-system/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type ConfirmHandler struct {
	repo repository.SubscriberRepository
}

func NewConfirmHandler(repo repository.SubscriberRepository) *ConfirmHandler {
	return &ConfirmHandler{repo: repo}
}

func (h *ConfirmHandler) Handle(c *fiber.Ctx) error {
	token := c.Query("token")

	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token is missing"})
	}

	sub, err := h.repo.FindByToken(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token not found"})
	}

	if sub == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "subscriber not found"})
	}

	if sub.Confirmed {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email already confirmed"})
	}

	if time.Now().After(sub.TokenExpiresAt) {
		return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "token got expired"})
	}

	err = h.repo.UpdateConfirmed(sub.ID, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "something went wrong"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "email confimred successfully"})
}
