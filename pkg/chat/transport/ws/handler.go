package ws

import (
	"app/pkg/chat/domain/entity"
	"app/pkg/chat/service/chat"
	"app/pkg/chat/service/chatroom"
	"app/pkg/chat/service/client"
	"app/pkg/database/redis"
	"app/pkg/exception"
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Handler handles WebSocket connections and events
type Handler struct {
	hub             *Hub
	clientService   client.ClientService
	chatService     chat.ChatService
	chatroomService chatroom.ChatroomService
}

// NewHandler creates a new WebSocket handler
func NewHandler(
	redisClient *redis.Client,
	clientService client.ClientService,
	chatService chat.ChatService,
	chatroomService chatroom.ChatroomService,
) *Handler {
	return &Handler{
		hub:             NewHub(redisClient, chatService, chatroomService),
		clientService:   clientService,
		chatService:     chatService,
		chatroomService: chatroomService,
	}
}

// RegisterRoutes registers WebSocket routes
func (h *Handler) RegisterRoutes(app fiber.Router) {
	ws := app.Group("/ws")

	// Middleware to upgrade HTTP connection to WebSocket
	ws.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket connection endpoint with authentication
	ws.Get("/", h.authenticateConnection(), websocket.New(h.handleConnection))
}

// authenticateConnection middleware authenticates the WebSocket connection
func (h *Handler) authenticateConnection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get client key from query params
		clientKey := c.Query("client_key")
		if clientKey == "" {
			return exception.Http(401, "Client key is required")
		}

		// Get auth token from query params
		token := c.Query("token")
		if token == "" {
			return exception.Http(401, "Auth token is required")
		}

		// Validate client key
		client, err := h.clientService.ValidateKey(c.Context(), clientKey)
		if err != nil {
			return err
		}

		// Authenticate user
		user, err := h.clientService.Authenticate(c.Context(), client.ID, token)
		if err != nil {
			return err
		}

		// Store authenticated data in context
		c.Locals("client", client)
		c.Locals("user", user)

		return c.Next()
	}
}

// handleConnection handles individual WebSocket connections
func (h *Handler) handleConnection(c *websocket.Conn) {
	// Get authenticated data from context
	user := c.Locals("user").(*entity.User)

	// Create new connection
	conn := &Connection{
		Socket: c,
		User:   user,
	}

	// Create new client
	client := NewClient(conn)

	// Register client with hub
	h.hub.register <- client

	// Start client message handlers
	go h.writePump(client)
	go h.readPump(client)
}

// writePump pumps messages from the hub to the WebSocket connection
func (h *Handler) writePump(client *Client) {
	for message := range client.Send {
		if err := client.Conn.Socket.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
	client.Conn.Socket.WriteMessage(websocket.CloseMessage, []byte{})
}

// readPump pumps messages from the WebSocket connection to the hub
func (h *Handler) readPump(client *Client) {
	defer func() {
		h.hub.unregister <- client
		client.Conn.Socket.Close()
	}()

	for {
		_, message, err := client.Conn.Socket.ReadMessage()
		if err != nil {
			break
		}

		// Parse the event
		var event Event
		if err := json.Unmarshal(message, &event); err != nil {
			continue
		}

		// Add user ID and timestamp to event
		event.UserID = client.Conn.User.ID
		event.Timestamp = time.Now().UnixMilli()

		// Handle different event types
		switch event.Type {
		case EventTypeMessage:
			h.handleChatMessage(client, &event)
		case EventTypeTypingStart, EventTypeTypingStop:
			h.handleTypingIndicator(client, &event)
		case EventTypeMessageRead:
			h.handleMessageRead(client, &event)
		}

		// Broadcast the event
		if data, err := json.Marshal(event); err == nil {
			h.hub.broadcast <- data
		}
	}
}

// handleChatMessage processes chat message events
func (h *Handler) handleChatMessage(client *Client, event *Event) {
	var payload MessagePayload
	if err := mapPayload(event.Payload, &payload); err != nil {
		return
	}

	// Create chat message using service
	params := chat.SendMessageParams{
		ChatroomID: event.ChatroomID,
		SenderID:   client.Conn.User.ID,
		Message:    payload.Message,
	}

	if msg, err := h.chatService.SendMessage(context.Background(), params); err == nil {
		event.Payload = msg
	}
}

// handleTypingIndicator processes typing indicator events
func (h *Handler) handleTypingIndicator(client *Client, event *Event) {
	var payload TypingPayload
	if err := mapPayload(event.Payload, &payload); err != nil {
		return
	}

	// Update payload with user info
	payload.UserID = client.Conn.User.ID
	payload.ChatroomID = event.ChatroomID
	event.Payload = payload
}

// handleMessageRead processes message read events
func (h *Handler) handleMessageRead(client *Client, event *Event) {
	var payload ReadPayload
	if err := mapPayload(event.Payload, &payload); err != nil {
		return
	}

	// Update payload with user info
	payload.UserID = client.Conn.User.ID
	event.Payload = payload
}

// mapPayload helper function to map interface{} to a specific type
func mapPayload(in interface{}, out interface{}) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
}

// TimeNow returns current Unix timestamp in milliseconds
func TimeNow() int64 {
	return time.Now().UnixMilli()
}
