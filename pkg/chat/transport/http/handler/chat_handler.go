package handler

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/domain/repository"
	"app/pkg/chat/service/chat"
	"app/pkg/chat/transport/http/dto"
	"app/pkg/chat/transport/http/middleware"
	"app/pkg/exception"
	"app/pkg/types/http"
	"app/pkg/types/pagination"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	chatService      chat.ChatService
	clientMiddleware *middleware.ClientMiddleware
	authMiddleware   *middleware.AuthMiddleware
}

func NewChatHandler(chatService chat.ChatService, clientMiddleware *middleware.ClientMiddleware, authMiddleware *middleware.AuthMiddleware) *ChatHandler {
	return &ChatHandler{
		chatService:      chatService,
		clientMiddleware: clientMiddleware,
		authMiddleware:   authMiddleware,
	}
}

// RegisterRoutes registers all routes for chat management
func (h *ChatHandler) RegisterRoutes(app fiber.Router) {
	v1 := app.Group("/v1")

	// Protected chat routes (requires client key and user authentication)
	chats := v1.Group("/chats", h.clientMiddleware.ValidateKey(), h.authMiddleware.Authenticate())

	// Chat message operations
	chats.Get("/", h.GetChats)                  // Get chat messages with filtering
	chats.Get("/:id", h.GetChat)                // Get single chat message
	chats.Put("/:id", h.UpdateChat)             // Update chat message
	chats.Delete("/:id", h.RemoveChat)          // Delete chat message
	chats.Post("/rooms/:roomId", h.SendMessage) // Send message to chatroom
	chats.Post("/direct", h.SendDirectMessage)  // Send direct message
}

// GetChats godoc
// @Summary Get chat messages
// @Description Retrieves chat messages with filtering and pagination
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param roomId query string false "Filter by chatroom ID"
// @Success 200 {object} http.GeneralResponse{data=http.PaginatedResponse{result=[]entity.Chat}}
// @Failure 401 {object} http.ErrorResponse
// @Router /v1/chats [get]
func (h *ChatHandler) GetChats(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	roomID := c.Query("roomId")

	pag := pagination.Pagination{
		Page:  page,
		Limit: limit,
	}

	filter := repository.ChatFilter{
		ChatroomID: roomID,
	}

	chats, total, err := h.chatService.GetChats(c.Context(), filter, pag)
	if err != nil {
		return err
	}

	metadata := pagination.Metadata{
		Pagination: pag,
		Total:      total,
		Count:      len(chats),
		HasPrev:    page > 1,
		HasNext:    len(chats) > 0 && int64(page*limit) < total,
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Chats fetched successfully",
		Data: map[string]interface{}{
			"metadata": metadata,
			"result":   chats,
		},
	})
}

// GetChat godoc
// @Summary Get a chat message
// @Description Retrieves a single chat message by ID
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chat ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Chat}
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/chats/{id} [get]
func (h *ChatHandler) GetChat(c *fiber.Ctx) error {
	id := c.Params("id")
	chat, err := h.chatService.GetChat(c.Context(), id)
	if err != nil {
		return err
	}
	if chat == nil {
		return exception.NotFound("Chat")
	}

	return c.JSON(http.GeneralResponse{
		Status: fiber.StatusOK,
		Data:   chat,
	})
}

// UpdateChat godoc
// @Summary Update a chat message
// @Description Updates an existing chat message
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chat ID"
// @Param chat body dto.UpdateChatRequest true "Chat details"
// @Success 200 {object} http.GeneralResponse{data=entity.Chat}
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/chats/{id} [put]
func (h *ChatHandler) UpdateChat(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if chat exists
	chat, err := h.chatService.GetChat(c.Context(), id)
	if err != nil {
		return err
	}
	if chat == nil {
		return exception.NotFound("Chat")
	}

	var req dto.UpdateChatRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	// Update only the message content
	chat.Message = req.Message

	if err := h.chatService.UpdateChat(c.Context(), chat); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Chat message updated successfully",
		Data:    chat,
	})
}

// RemoveChat godoc
// @Summary Delete a chat message
// @Description Deletes an existing chat message
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chat ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Chat}
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/chats/{id} [delete]
func (h *ChatHandler) RemoveChat(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if chat exists
	chat, err := h.chatService.GetChat(c.Context(), id)
	if err != nil {
		return err
	}
	if chat == nil {
		return exception.NotFound("Chat")
	}

	if err := h.chatService.RemoveChat(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Chat message deleted successfully",
	})
}

// SendMessage godoc
// @Summary Send a message to a chatroom
// @Description Sends a new message to a specific chatroom
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param roomId path string true "Chatroom ID"
// @Param message body dto.SendMessageRequest true "Message details"
// @Success 201 {object} http.GeneralResponse{data=entity.Chat}
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/chats/rooms/{roomId} [post]
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	roomID := c.Params("roomId")
	user := c.Locals("user").(*entity.User)

	var req dto.SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	params := chat.SendMessageParams{
		ChatroomID: roomID,
		SenderID:   user.ID,
		Message:    req.Message,
	}

	chat, err := h.chatService.SendMessage(c.Context(), params)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Message sent successfully",
		Data:    chat,
	})
}

// SendDirectMessage godoc
// @Summary Send a direct message
// @Description Sends a direct message to another user
// @Tags chats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param message body dto.SendDirectMessageRequest true "Message details"
// @Success 201 {object} http.GeneralResponse{data=entity.Chat}
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/chats/direct [post]
func (h *ChatHandler) SendDirectMessage(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.User)

	var req dto.SendDirectMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	params := chat.SendDirectMessageParams{
		SenderID:   user.ID,
		ReceiverID: req.ReceiverID,
		Message:    req.Message,
	}

	chat, chatroom, err := h.chatService.SendDirectMessage(c.Context(), params)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Direct message sent successfully",
		Data: fiber.Map{
			"chat":     chat,
			"chatroom": chatroom,
		},
	})
}
