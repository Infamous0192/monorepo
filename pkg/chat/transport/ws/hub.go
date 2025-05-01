package ws

import (
	"app/pkg/chat/service/chat"
	"app/pkg/chat/service/chatroom"
	"app/pkg/database/redis"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	// Redis key prefixes
	connectedUsersKey = "ws:connected_users"
	userSocketKey     = "ws:user_socket:%s" // Format with user ID
	socketUserKey     = "ws:socket_user:%s" // Format with socket ID

	// Redis expiration times
	socketExpiration = 24 * time.Hour
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Redis client for managing connected users
	redisClient *redis.Client

	// Chat service for handling messages
	chatService chat.ChatService

	// Chatroom service for managing rooms
	chatroomService chatroom.ChatroomService

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub(redisClient *redis.Client, chatService chat.ChatService, chatroomService chatroom.ChatroomService) *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		broadcast:       make(chan []byte),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		redisClient:     redisClient,
		chatService:     chatService,
		chatroomService: chatroomService,
	}
}

// Run starts the hub's main event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)

		case client := <-h.unregister:
			h.handleUnregister(client)

		case message := <-h.broadcast:
			h.handleBroadcast(message)
		}
	}
}

// handleRegister processes a new client registration
func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	h.clients[client] = true
	h.mu.Unlock()

	// Store user connection in Redis
	if err := h.storeUserConnection(client); err != nil {
		fmt.Printf("Error storing user connection: %v\n", err)
		return
	}

	// Notify other clients about the new user
	h.broadcastUserStatus(client, true)
}

// handleUnregister processes a client disconnection
func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)
	}
	h.mu.Unlock()

	// Remove user connection from Redis
	if err := h.removeUserConnection(client); err != nil {
		fmt.Printf("Error removing user connection: %v\n", err)
	}

	// Notify other clients about the user leaving
	h.broadcastUserStatus(client, false)
}

// handleBroadcast processes and broadcasts a message to relevant clients
func (h *Hub) handleBroadcast(message []byte) {
	var event Event
	if err := json.Unmarshal(message, &event); err != nil {
		fmt.Printf("Error unmarshaling event: %v\n", err)
		return
	}

	// If the message is for a specific chatroom, only send to clients in that room
	if event.ChatroomID != "" {
		h.broadcastToChatroom(event.ChatroomID, message)
		return
	}

	// Otherwise, broadcast to all clients
	h.broadcastToAll(message)
}

// broadcastToChatroom sends a message to all clients in a specific chatroom
func (h *Hub) broadcastToChatroom(chatroomID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.Chatrooms[chatroomID] {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// broadcastToAll sends a message to all connected clients
func (h *Hub) broadcastToAll(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

// broadcastUserStatus notifies clients about a user's connection status
func (h *Hub) broadcastUserStatus(client *Client, isConnected bool) {
	eventType := EventTypeConnect
	if !isConnected {
		eventType = EventTypeDisconnect
	}

	event := Event{
		Type:      eventType,
		UserID:    client.Conn.User.ID,
		Timestamp: time.Now().UnixMilli(),
		Payload: map[string]interface{}{
			"user": client.Conn.User,
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Error marshaling user status event: %v\n", err)
		return
	}

	h.broadcastToAll(data)
}

// storeUserConnection stores user connection details in Redis
func (h *Hub) storeUserConnection(client *Client) error {
	ctx := context.Background()
	userID := client.Conn.User.ID
	socketID := client.Conn.Socket.LocalAddr().String()

	// Store user's online status
	if err := h.redisClient.SAdd(ctx, connectedUsersKey, userID); err != nil {
		return fmt.Errorf("error storing user online status: %v", err)
	}

	// Store user-socket mapping
	if err := h.redisClient.Set(ctx, fmt.Sprintf(userSocketKey, userID), socketID, socketExpiration); err != nil {
		return fmt.Errorf("error storing user-socket mapping: %v", err)
	}

	// Store socket-user mapping
	if err := h.redisClient.Set(ctx, fmt.Sprintf(socketUserKey, socketID), userID, socketExpiration); err != nil {
		return fmt.Errorf("error storing socket-user mapping: %v", err)
	}

	return nil
}

// removeUserConnection removes user connection details from Redis
func (h *Hub) removeUserConnection(client *Client) error {
	ctx := context.Background()
	userID := client.Conn.User.ID
	socketID := client.Conn.Socket.LocalAddr().String()

	// Remove user's online status
	if err := h.redisClient.SRem(ctx, connectedUsersKey, userID); err != nil {
		return fmt.Errorf("error removing user online status: %v", err)
	}

	// Remove user-socket mapping
	if err := h.redisClient.Del(ctx, fmt.Sprintf(userSocketKey, userID)); err != nil {
		return fmt.Errorf("error removing user-socket mapping: %v", err)
	}

	// Remove socket-user mapping
	if err := h.redisClient.Del(ctx, fmt.Sprintf(socketUserKey, socketID)); err != nil {
		return fmt.Errorf("error removing socket-user mapping: %v", err)
	}

	return nil
}

// IsUserConnected checks if a user is currently connected
func (h *Hub) IsUserConnected(userID string) (bool, error) {
	return h.redisClient.SIsMember(context.Background(), connectedUsersKey, userID)
}

// GetConnectedUsers returns a list of all connected user IDs
func (h *Hub) GetConnectedUsers() ([]string, error) {
	return h.redisClient.SMembers(context.Background(), connectedUsersKey)
}

// GetUserSocket returns the socket ID for a connected user
func (h *Hub) GetUserSocket(userID string) (string, error) {
	return h.redisClient.Get(context.Background(), fmt.Sprintf(userSocketKey, userID))
}

// GetSocketUser returns the user ID associated with a socket
func (h *Hub) GetSocketUser(socketID string) (string, error) {
	return h.redisClient.Get(context.Background(), fmt.Sprintf(socketUserKey, socketID))
}
