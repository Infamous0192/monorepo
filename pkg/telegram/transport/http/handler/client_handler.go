package handler

import (
	"app/pkg/exception"
	"app/pkg/middleware"
	"app/pkg/telegram/domain/entity"
	"app/pkg/telegram/service/client"
	"app/pkg/telegram/transport/http/dto"
	httpMiddleware "app/pkg/telegram/transport/http/middleware"
	"app/pkg/types/http"
	"app/pkg/types/pagination"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ClientHandler handles HTTP requests for Telegram bot client management
type ClientHandler struct {
	clientService    client.ClientService
	keyMiddleware    *middleware.KeyMiddleware
	clientMiddleware *httpMiddleware.ClientMiddleware
}

// NewClientHandler creates a new instance of ClientHandler
func NewClientHandler(clientService client.ClientService, keyMiddleware *middleware.KeyMiddleware, clientMiddleware *httpMiddleware.ClientMiddleware) *ClientHandler {
	return &ClientHandler{
		clientService:    clientService,
		keyMiddleware:    keyMiddleware,
		clientMiddleware: clientMiddleware,
	}
}

// RegisterRoutes registers all routes for client management
func (h *ClientHandler) RegisterRoutes(app fiber.Router) {
	v1 := app.Group("/v1")

	// Protected client routes (requires API key)
	clients := v1.Group("/clients", h.keyMiddleware.ValidateKey())

	// Client operations
	clients.Get("/", h.GetClients)         // Get all clients with pagination
	clients.Get("/:id", h.GetClient)       // Get single client
	clients.Post("/", h.CreateClient)      // Create new client
	clients.Put("/:id", h.UpdateClient)    // Update client
	clients.Delete("/:id", h.DeleteClient) // Delete client

	// Client-specific routes (requires client key)
	clientAPI := v1.Group("/client", h.clientMiddleware.ValidateKey())
	clientAPI.Get("/profile", h.GetClientProfile) // Get client profile
}

// GetClients godoc
// @Summary Get all clients
// @Description Retrieves all Telegram bot clients with pagination
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} http.GeneralResponse{data=pagination.PaginatedResult{result=[]entity.Client}}
// @Failure 401 {object} http.ErrorResponse
// @Router /v1/clients [get]
func (h *ClientHandler) GetClients(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pag := pagination.Pagination{
		Page:  page,
		Limit: limit,
	}

	clients, total, err := h.clientService.GetAll(c.Context(), pag)
	if err != nil {
		return err
	}

	metadata := pagination.Metadata{
		Pagination: pag,
		Total:      total,
		Count:      len(clients),
		HasPrev:    page > 1,
		HasNext:    len(clients) > 0 && int64(page*limit) < total,
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Clients fetched successfully",
		Data: map[string]interface{}{
			"metadata": metadata,
			"result":   clients,
		},
	})
}

// GetClient godoc
// @Summary Get a client
// @Description Retrieves a single Telegram bot client by ID
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Client ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/clients/{id} [get]
func (h *ClientHandler) GetClient(c *fiber.Ctx) error {
	id := c.Params("id")
	client, err := h.clientService.Get(c.Context(), id)
	if err != nil {
		return err
	}
	if client == nil {
		return exception.NotFound("Client")
	}

	return c.JSON(http.GeneralResponse{
		Status: fiber.StatusOK,
		Data:   client,
	})
}

// CreateClient godoc
// @Summary Create a client
// @Description Creates a new Telegram bot client
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param client body dto.CreateClientRequest true "Client details"
// @Success 201 {object} http.GeneralResponse
// @Failure 400 {object} http.ErrorResponse
// @Router /v1/clients [post]
func (h *ClientHandler) CreateClient(c *fiber.Ctx) error {
	var req dto.CreateClientRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	// Create client entity
	client := &entity.Client{
		Token:            req.Token,
		Username:         req.Username,
		Name:             req.Name,
		Description:      req.Description,
		BotType:          req.BotType,
		WebhookURL:       req.WebhookURL,
		MaxConnections:   req.MaxConnections,
		AllowedUpdates:   req.AllowedUpdates,
		CreatedTimestamp: time.Now().UnixMilli(),
		UpdatedTimestamp: time.Now().UnixMilli(),
	}

	// Set default status if not provided
	if req.Status == "" {
		client.Status = "active"
	} else {
		client.Status = req.Status
	}

	if err := h.clientService.Create(c.Context(), client); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Client created successfully",
		Data:    client,
	})
}

// UpdateClient godoc
// @Summary Update a client
// @Description Updates an existing Telegram bot client
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Client ID"
// @Param client body dto.UpdateClientRequest true "Client details"
// @Success 200 {object} http.GeneralResponse
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/clients/{id} [put]
func (h *ClientHandler) UpdateClient(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if client exists
	existingClient, err := h.clientService.Get(c.Context(), id)
	if err != nil {
		return err
	}
	if existingClient == nil {
		return exception.NotFound("Client")
	}

	var req dto.UpdateClientRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	// Update client fields
	existingClient.Token = req.Token
	existingClient.Username = req.Username
	existingClient.Name = req.Name
	existingClient.Description = req.Description
	existingClient.BotType = req.BotType
	existingClient.WebhookURL = req.WebhookURL
	existingClient.Status = req.Status
	existingClient.MaxConnections = req.MaxConnections
	existingClient.AllowedUpdates = req.AllowedUpdates
	existingClient.UpdatedTimestamp = time.Now().UnixMilli()

	if err := h.clientService.Update(c.Context(), existingClient); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Client updated successfully",
		Data:    existingClient,
	})
}

// DeleteClient godoc
// @Summary Delete a client
// @Description Deletes an existing Telegram bot client
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Client ID"
// @Success 200 {object} http.GeneralResponse
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/clients/{id} [delete]
func (h *ClientHandler) DeleteClient(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if client exists
	client, err := h.clientService.Get(c.Context(), id)
	if err != nil {
		return err
	}
	if client == nil {
		return exception.NotFound("Client")
	}

	if err := h.clientService.Delete(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Client deleted successfully",
	})
}

// GetClientProfile godoc
// @Summary Get client profile
// @Description Retrieves the profile of the authenticated client
// @Tags client
// @Accept json
// @Produce json
// @Security ClientKeyAuth
// @Success 200 {object} http.GeneralResponse
// @Failure 401 {object} http.ErrorResponse
// @Router /v1/client/profile [get]
func (h *ClientHandler) GetClientProfile(c *fiber.Ctx) error {
	// Get client from context (set by middleware)
	client := c.Locals("client").(*entity.Client)

	return c.JSON(http.GeneralResponse{
		Status: fiber.StatusOK,
		Data:   client,
	})
}
