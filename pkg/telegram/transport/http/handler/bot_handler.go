package handler

import (
	"app/pkg/exception"
	"app/pkg/telegram/service/bot"
	"app/pkg/telegram/service/client"
	"app/pkg/types/http"

	"github.com/gofiber/fiber/v2"
)

type BotHandler struct {
	botService    *bot.BotService
	clientService client.ClientService
}

// NewBotHandler creates a new instance of BotHandler
func NewBotHandler(botService *bot.BotService, clientService client.ClientService) *BotHandler {
	return &BotHandler{
		botService:    botService,
		clientService: clientService,
	}
}

// RegisterRoutes registers all bot-related routes
func (h *BotHandler) RegisterRoutes(app fiber.Router) {
	botGroup := app.Group("/v1/bots")

	botGroup.Post("/start", h.StartBot)
	botGroup.Post("/stop/:clientId", h.StopBot)
	botGroup.Post("/stop-all", h.StopAllBots)
	botGroup.Get("/status/:clientId", h.GetBotStatus)
	botGroup.Get("/status", h.GetAllBotsStatus)
}

// StartBot godoc
// @Summary Start a Telegram bot
// @Description Start a Telegram bot for a given client
// @Tags bots
// @Accept json
// @Produce json
// @Param clientId body string true "Client ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} http.ErrorResponse
// @Router /v1/bots/start [post]
func (h *BotHandler) StartBot(c *fiber.Ctx) error {
	var req struct {
		ClientID string `json:"clientId"`
	}

	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	// Get client details
	client, err := h.clientService.Get(c.Context(), req.ClientID)
	if err != nil {
		return err
	}
	if client == nil {
		return exception.NotFound("Client")
	}

	// Start the bot
	if err := h.botService.StartBot(c.Context(), client); err != nil {
		return exception.InternalError(err.Error())
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Bot started successfully",
		Data: fiber.Map{
			"clientId": client.ID,
		},
	})
}

// StopBot godoc
// @Summary Stop a Telegram bot
// @Description Stop a Telegram bot for a given client
// @Tags bots
// @Accept json
// @Produce json
// @Param clientId path string true "Client ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 400 {object} http.ErrorResponse
// @Router /v1/bots/stop/{clientId} [post]
func (h *BotHandler) StopBot(c *fiber.Ctx) error {
	clientID := c.Params("clientId")

	client, err := h.clientService.Get(c.Context(), clientID)
	if err != nil {
		return err
	}

	if err := h.botService.StopBot(c.Context(), client); err != nil {
		return exception.InternalError(err.Error())
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Bot stopped successfully",
		Data: fiber.Map{
			"clientId": clientID,
		},
	})
}

// StopAllBots godoc
// @Summary Stop all Telegram bots
// @Description Stop all running Telegram bots
// @Tags bots
// @Accept json
// @Produce json
// @Success 200 {object} http.GeneralResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /v1/bots/stop-all [post]
func (h *BotHandler) StopAllBots(c *fiber.Ctx) error {
	if err := h.botService.StopAllBots(c.Context()); err != nil {
		return exception.InternalError(err.Error())
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "All bots stopped successfully",
	})
}

// GetBotStatus godoc
// @Summary Get bot status
// @Description Get the status of a specific Telegram bot
// @Tags bots
// @Accept json
// @Produce json
// @Param clientId path string true "Client ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/bots/status/{clientId} [get]
func (h *BotHandler) GetBotStatus(c *fiber.Ctx) error {
	clientID := c.Params("clientId")

	_, err := h.botService.GetBot(clientID)
	if err != nil {
		return exception.NotFound("Bot")
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Bot status retrieved successfully",
		Data: fiber.Map{
			"clientId": clientID,
			"status":   "running",
		},
	})
}

// GetAllBotsStatus godoc
// @Summary Get all bots status
// @Description Get the status of all running Telegram bots
// @Tags bots
// @Accept json
// @Produce json
// @Success 200 {object} http.GeneralResponse
// @Router /v1/bots/status [get]
func (h *BotHandler) GetAllBotsStatus(c *fiber.Ctx) error {
	bots := h.botService.GetAllBots()
	status := make([]fiber.Map, 0, len(bots))

	for clientID := range bots {
		status = append(status, fiber.Map{
			"clientId": clientID,
			"status":   "running",
		})
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "All bot statuses retrieved successfully",
		Data:    status,
	})
}
