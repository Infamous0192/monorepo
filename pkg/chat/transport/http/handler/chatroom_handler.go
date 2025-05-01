package handler

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/domain/repository"
	"app/pkg/chat/service/chatroom"
	"app/pkg/chat/transport/http/dto"
	"app/pkg/chat/transport/http/middleware"
	"app/pkg/exception"
	"app/pkg/types/http"
	"app/pkg/types/pagination"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ChatroomHandler struct {
	chatroomService  chatroom.ChatroomService
	clientMiddleware *middleware.ClientMiddleware
	authMiddleware   *middleware.AuthMiddleware
}

func NewChatroomHandler(chatroomService chatroom.ChatroomService, clientMiddleware *middleware.ClientMiddleware, authMiddleware *middleware.AuthMiddleware) *ChatroomHandler {
	return &ChatroomHandler{
		chatroomService:  chatroomService,
		clientMiddleware: clientMiddleware,
		authMiddleware:   authMiddleware,
	}
}

// RegisterRoutes registers all routes for chatroom management
func (h *ChatroomHandler) RegisterRoutes(app fiber.Router) {
	v1 := app.Group("/v1")

	// Protected chatroom routes (requires client key and user authentication)
	chatrooms := v1.Group("/chatrooms", h.clientMiddleware.ValidateKey(), h.authMiddleware.Authenticate())

	// Chatroom operations
	chatrooms.Post("/", h.CreateChatroom)      // Create new chatroom
	chatrooms.Get("/", h.GetChatrooms)         // List chatrooms
	chatrooms.Get("/:id", h.GetChatroom)       // Get single chatroom
	chatrooms.Put("/:id", h.UpdateChatroom)    // Update chatroom
	chatrooms.Delete("/:id", h.DeleteChatroom) // Delete chatroom

	// Participant management
	chatrooms.Post("/:id/participants", h.AddParticipant)                    // Add participant
	chatrooms.Delete("/:id/participants/:userId", h.RemoveParticipant)       // Remove participant
	chatrooms.Put("/:id/participants/:userId/role", h.UpdateParticipantRole) // Update participant role
	chatrooms.Post("/:id/participants/:userId/mute", h.MuteParticipant)      // Mute participant
	chatrooms.Post("/:id/participants/:userId/unmute", h.UnmuteParticipant)  // Unmute participant
}

// CreateChatroom godoc
// @Summary Create a new chatroom
// @Description Creates a new chatroom with the provided details
// @Tags chatrooms
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param chatroom body dto.CreateChatroomRequest true "Chatroom details"
// @Success 201 {object} http.GeneralResponse{data=entity.Chatroom}
// @Failure 400 {object} http.ErrorResponse
// @Router /v1/chatrooms [post]
func (h *ChatroomHandler) CreateChatroom(c *fiber.Ctx) error {
	user := c.Locals("user").(*entity.User)

	var req dto.CreateChatroomRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	params := chatroom.CreateChatroomParams{
		Name:         req.Name,
		Type:         entity.ChatroomType(req.Type),
		IsGroup:      req.IsGroup,
		Creator:      user.ID,
		Participants: req.Participants,
	}

	chatroom, err := h.chatroomService.CreateChatroom(c.Context(), params)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(http.GeneralResponse{
		Status:  fiber.StatusCreated,
		Message: "Chatroom created successfully",
		Data:    chatroom,
	})
}

// GetChatrooms godoc
// @Summary Get chatrooms
// @Description Retrieves chatrooms with filtering and pagination
// @Tags chatrooms
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} http.GeneralResponse{data=http.PaginatedResponse{result=[]entity.Chatroom}}
// @Failure 401 {object} http.ErrorResponse
// @Router /v1/chatrooms [get]
func (h *ChatroomHandler) GetChatrooms(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pag := pagination.Pagination{
		Page:  page,
		Limit: limit,
	}

	// Create empty filter for now, add filter params as needed
	filter := repository.ChatroomFilter{}

	chatrooms, total, err := h.chatroomService.GetChatrooms(c.Context(), filter, pag)
	if err != nil {
		return err
	}

	metadata := pagination.Metadata{
		Pagination: pag,
		Total:      total,
		Count:      len(chatrooms),
		HasPrev:    page > 1,
		HasNext:    len(chatrooms) > 0 && int64(page*limit) < total,
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Chatrooms fetched successfully",
		Data: map[string]interface{}{
			"metadata": metadata,
			"result":   chatrooms,
		},
	})
}

// GetChatroom godoc
// @Summary Get a chatroom
// @Description Retrieves a single chatroom by ID
// @Tags chatrooms
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chatroom ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Chatroom}
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/chatrooms/{id} [get]
func (h *ChatroomHandler) GetChatroom(c *fiber.Ctx) error {
	id := c.Params("id")
	chatroom, err := h.chatroomService.GetChatroom(c.Context(), id)
	if err != nil {
		return err
	}
	if chatroom == nil {
		return exception.NotFound("Chatroom")
	}

	return c.JSON(http.GeneralResponse{
		Status: fiber.StatusOK,
		Data:   chatroom,
	})
}

// UpdateChatroom godoc
// @Summary Update a chatroom
// @Description Updates an existing chatroom
// @Tags chatrooms
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chatroom ID"
// @Param chatroom body dto.UpdateChatroomRequest true "Chatroom details"
// @Success 200 {object} http.GeneralResponse{data=entity.Chatroom}
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/chatrooms/{id} [put]
func (h *ChatroomHandler) UpdateChatroom(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if chatroom exists
	existingChatroom, err := h.chatroomService.GetChatroom(c.Context(), id)
	if err != nil {
		return err
	}
	if existingChatroom == nil {
		return exception.NotFound("Chatroom")
	}

	var req dto.UpdateChatroomRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	chatroomType := entity.ChatroomType(req.Type)
	params := chatroom.UpdateChatroomParams{
		ID:   id,
		Name: &req.Name,
		Type: &chatroomType,
	}

	updatedChatroom, err := h.chatroomService.UpdateChatroom(c.Context(), params)
	if err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Chatroom updated successfully",
		Data:    updatedChatroom,
	})
}

// DeleteChatroom godoc
// @Summary Delete a chatroom
// @Description Deletes an existing chatroom
// @Tags chatrooms
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chatroom ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Chatroom}
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/chatrooms/{id} [delete]
func (h *ChatroomHandler) DeleteChatroom(c *fiber.Ctx) error {
	id := c.Params("id")

	// Check if chatroom exists
	chatroom, err := h.chatroomService.GetChatroom(c.Context(), id)
	if err != nil {
		return err
	}
	if chatroom == nil {
		return exception.NotFound("Chatroom")
	}

	if err := h.chatroomService.DeleteChatroom(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Chatroom deleted successfully",
	})
}

// AddParticipant godoc
// @Summary Add a participant to a chatroom
// @Description Adds a new participant to an existing chatroom
// @Tags chatrooms
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chatroom ID"
// @Param participant body dto.AddParticipantRequest true "Participant details"
// @Success 200 {object} http.GeneralResponse{data=entity.Chatroom}
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/chatrooms/{id}/participants [post]
func (h *ChatroomHandler) AddParticipant(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.AddParticipantRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	if err := h.chatroomService.AddParticipant(c.Context(), id, req.UserID); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Participant added successfully",
	})
}

// RemoveParticipant godoc
// @Summary Remove a participant from a chatroom
// @Description Removes a participant from an existing chatroom
// @Tags chatrooms
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chatroom ID"
// @Param userId path string true "User ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Chatroom}
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/chatrooms/{id}/participants/{userId} [delete]
func (h *ChatroomHandler) RemoveParticipant(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Params("userId")

	if err := h.chatroomService.RemoveParticipant(c.Context(), id, userID); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Participant removed successfully",
	})
}

// UpdateParticipantRole godoc
// @Summary Update a participant's role
// @Description Updates the role of a participant in a chatroom
// @Tags chatrooms
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chatroom ID"
// @Param userId path string true "User ID"
// @Param role body dto.UpdateParticipantRequest true "Role details"
// @Success 200 {object} http.GeneralResponse{data=entity.Chatroom}
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/chatrooms/{id}/participants/{userId}/role [put]
func (h *ChatroomHandler) UpdateParticipantRole(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Params("userId")

	var req dto.UpdateParticipantRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	params := chatroom.UpdateParticipantRoleParams{
		ChatroomID:    id,
		ParticipantID: userID,
		NewRole:       entity.ParticipantRole(req.Role),
	}

	if err := h.chatroomService.UpdateParticipantRole(c.Context(), params); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Participant role updated successfully",
	})
}

// MuteParticipant godoc
// @Summary Mute a participant
// @Description Mutes a participant in a chatroom for a specified duration
// @Tags chatrooms
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chatroom ID"
// @Param userId path string true "User ID"
// @Param mute body dto.MuteParticipantRequest true "Mute details"
// @Success 200 {object} http.GeneralResponse{data=entity.Chatroom}
// @Failure 400,404 {object} http.ErrorResponse
// @Router /v1/chatrooms/{id}/participants/{userId}/mute [post]
func (h *ChatroomHandler) MuteParticipant(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Params("userId")

	var req dto.MuteParticipantRequest
	if err := c.BodyParser(&req); err != nil {
		return exception.BadRequest("Invalid request body")
	}

	params := chatroom.MuteParticipantParams{
		ChatroomID:    id,
		ParticipantID: userID,
		Duration:      time.Duration(req.Duration) * time.Minute,
	}

	if err := h.chatroomService.MuteParticipant(c.Context(), params); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Participant muted successfully",
	})
}

// UnmuteParticipant godoc
// @Summary Unmute a participant
// @Description Removes the mute from a participant in a chatroom
// @Tags chatrooms
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path string true "Chatroom ID"
// @Param userId path string true "User ID"
// @Success 200 {object} http.GeneralResponse{data=entity.Chatroom}
// @Failure 404 {object} http.ErrorResponse
// @Router /v1/chatrooms/{id}/participants/{userId}/unmute [post]
func (h *ChatroomHandler) UnmuteParticipant(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Params("userId")

	if err := h.chatroomService.UnmuteParticipant(c.Context(), id, userID); err != nil {
		return err
	}

	return c.JSON(http.GeneralResponse{
		Status:  fiber.StatusOK,
		Message: "Participant unmuted successfully",
	})
}
