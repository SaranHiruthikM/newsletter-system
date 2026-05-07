package handlers

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *fiber.Ctx) error {
	err := h.db.Ping()

	if err != nil {
		return c.Status(503).JSON(fiber.Map{"status": "unhealthy"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "healthy"})
}
