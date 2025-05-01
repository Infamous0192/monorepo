package handler

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/service/client"
	"app/pkg/chat/transport/http/dto"
	"app/pkg/exception"
	"app/pkg/middleware"
	"app/pkg/types/http"
	"app/pkg/types/pagination"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ClientHandler struct {
	clientService client.ClientService
	keyMiddleware *middleware.KeyMiddleware
}

func NewClientHandler(clientService client.ClientService, keyMiddleware *middleware.KeyMiddleware) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
		keyMiddleware: keyMiddleware,
	}
}

// RegisterRoutes registers all routes for client management
func (h *ClientHandler) RegisterRoutes(app fiber.Router) {
	v1 := app.Group("/v1")

	// Public routes
	clients := v1.Group("/clients")
	clients.Get("/validate", h.ValidateClientKey) // For testing client keys

	// Admin protected routes
	adminClients := v1.Group("/admin/clients", h.keyMiddleware.ValidateKey())
	adminClients.Post("/", h.CreateClient)
	adminClients.Get("/", h.GetClients)
	adminClients.Get("/:id", h.GetClient)
	adminClients.Put("/:id", h.UpdateClient)
	adminClients.Delete("/:id", h.DeleteClient)
}

// ValidateClientKey godoc
// @Summary Validate a client key
// @Description Validates a client key and returns the client if valid
// @Tags clients
// @Accept json
// @Produce json
// @Param X-Client-Key header string true "Client Key"
// @Success 200 {object} http.GeneralResponse{data=entity.Client}
// @Failure 401 {object} http.ErrorResponse
// @Router /v1/clients/validate [get]
func (h *ClientHandler) ValidateClientKey(c *fiber.Ctx) error {
	clientKey := c.Get("X-Client-Key")
	client, err := h.clientService.ValidateKey(c.Context(), clientKey)
	if err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status: fiber.StatusOK,
		Data:   client,
	})
}

// CreateClient godoc
// @Summary Create a new client
// @Description Creates a new client with the provided details
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param client body dto.CreateClientRequest true "Client details"
// @Success 201 {object} http.GeneralResponse{data=entity.Client}
// @Failure 400 {object} http.ErrorResponse
// @Router /v1/admin/clients [post]
func (h *ClientHandler) CreateClient(c *fiber.Ctx) error {
	var req dto.CreateClientRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	client := &entity.Client{
		Name:         req.Name,
		Description:  req.Description,
		ClientKey:    req.ClientKey,
		AuthEndpoint: req.AuthEndpoint,
		Status:       "active", // Default status for new clients
	}

	if err := h.clientService.CreateClient(c.Context(), client); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Client created successfully",
		Data:    client,
	})
}

// GetClients godoc
// @Summary Get clients
// @Description Retrieves clients with filtering and pagination
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} http.GeneralResponse{data=http.PaginatedResponse{result=[]entity.Client}}
// @Failure 401 {object} http.ErrorResponse
// @Router /v1/clients [get]
func (h *ClientHandler) GetClients(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pag := pagination.Pagination{
		Page:  page,
		Limit: limit,
	}

	clients, total, err := h.clientService.GetClients(c.Context(), pag)
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
// @Summary Get a client by ID
// @Description Retrieves a client by its ID
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Client ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Client}
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/admin/clients/{id} [get]
func (h *ClientHandler) GetClient(c *fiber.Ctx) error {
	id := c.Params("id")
	client, err := h.clientService.GetClient(c.Context(), id)
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

// UpdateClient godoc
// @Summary Update a client
// @Description Updates an existing client
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Client ID"
// @Param client body dto.UpdateClientRequest true "Client details"
// @Success 200 {object} http.GeneralResponse{data=entity.Client}
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/admin/clients/{id} [put]
func (h *ClientHandler) UpdateClient(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if client exists
	existingClient, err := h.clientService.GetClient(c.Context(), id)
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

	client := &entity.Client{
		ID:           id,
		Name:         req.Name,
		Description:  req.Description,
		ClientKey:    req.ClientKey,
		AuthEndpoint: req.AuthEndpoint,
		Status:       req.Status,
	}

	if err := h.clientService.UpdateClient(c.Context(), client); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Client updated successfully",
		Data:    client,
	})
}

// DeleteClient godoc
// @Summary Delete a client
// @Description Deletes an existing client
// @Tags clients
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Client ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Client}
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/admin/clients/{id} [delete]
func (h *ClientHandler) DeleteClient(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if client exists
	client, err := h.clientService.GetClient(c.Context(), id)
	if err != nil {
		return err
	}
	if client == nil {
		return exception.NotFound("Client")
	}

	if err := h.clientService.DeleteClient(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Client deleted successfully",
	})
}
