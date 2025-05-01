package handler

import (
	"app/pkg/exception"
	"app/pkg/telegram/service/bot"
	"app/pkg/telegram/service/client"
	"app/pkg/types/http"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/telebot.v4"
)

type WebhookHandler struct {
	botService    *bot.BotService
	clientService client.ClientService
}

// NewWebhookHandler creates a new instance of WebhookHandler
func NewWebhookHandler(botService *bot.BotService, clientService client.ClientService) *WebhookHandler {
	return &WebhookHandler{
		botService:    botService,
		clientService: clientService,
	}
}

// RegisterRoutes registers all webhook-related routes
func (h *WebhookHandler) RegisterRoutes(app fiber.Router) {
	webhookGroup := app.Group("/v1/webhook")
	webhookGroup.Post("/:clientId", h.HandleUpdate)
}

// HandleUpdate godoc
// @Summary Handle Telegram webhook update
// @Description Process incoming updates from Telegram webhook
// @Tags webhook
// @Accept json
// @Produce json
// @Param clientId path string true "Client ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/webhook/{clientId} [post]
func (h *WebhookHandler) HandleUpdate(c *fiber.Ctx) error {
	clientID := c.Params("clientId")

	// Get the bot instance
	_, err := h.botService.GetBot(clientID)
	if err != nil {
		return exception.NotFound("Bot not found")
	}

	// Parse the update
	var update telebot.Update
	if err := json.Unmarshal(c.Body(), &update); err != nil {
		return exception.BadRequest("Invalid update format")
	}

	// Process the update
	// if err := bot.HandleUpdate(c.Context(), &update); err != nil {
	// 	return exception.InternalError("Failed to process update")
	// }

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Update processed successfully",
	})
}
