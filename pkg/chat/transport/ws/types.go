package ws

import (
	"app/pkg/chat/domain/entity"
	"encoding/json"

	"github.com/gofiber/websocket/v2"
)

// EventType represents different types of WebSocket events
type EventType string

const (
	// Connection events
	EventTypeConnect    EventType = "connect"
	EventTypeDisconnect EventType = "disconnect"

	// Chat events
	EventTypeMessage      EventType = "message"
	EventTypeMessageRead  EventType = "message_read"
	EventTypeTypingStart  EventType = "typing_start"
	EventTypeTypingStop   EventType = "typing_stop"
	EventTypeUserJoin     EventType = "user_join"
	EventTypeUserLeave    EventType = "user_leave"
	EventTypeUserMuted    EventType = "user_muted"
	EventTypeUserUnmuted  EventType = "user_unmuted"
	EventTypeRoleUpdated  EventType = "role_updated"
	EventTypeChatroomMeta EventType = "chatroom_meta"
	EventTypeError        EventType = "error"
)

// Event represents a WebSocket event
type Event struct {
	Type       EventType   `json:"type"`
	ChatroomID string      `json:"chatroomId,omitempty"`
	UserID     string      `json:"userId,omitempty"`
	Payload    interface{} `json:"payload,omitempty"`
	Timestamp  int64       `json:"timestamp"`
}

// Connection represents a WebSocket connection with user information
type Connection struct {
	Socket *websocket.Conn
	User   *entity.User
}

// Client represents a connected WebSocket client
type Client struct {
	Conn      *Connection
	Send      chan []byte
	Chatrooms map[string]bool // Map of chatroom IDs the client is subscribed to
}

// MessagePayload represents a chat message event payload
type MessagePayload struct {
	ID       string                 `json:"id,omitempty"`
	Message  string                 `json:"message"`
	Type     string                 `json:"type"`
	ReplyTo  string                 `json:"replyTo,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// TypingPayload represents a typing indicator event payload
type TypingPayload struct {
	ChatroomID string `json:"chatroomId"`
	UserID     string `json:"userId"`
	Status     bool   `json:"status"` // true for typing, false for stopped typing
}

// ReadPayload represents a message read event payload
type ReadPayload struct {
	ChatroomID string `json:"chatroomId"`
	MessageID  string `json:"messageId"`
	UserID     string `json:"userId"`
}

// ErrorPayload represents an error event payload
type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewClient creates a new WebSocket client
func NewClient(conn *Connection) *Client {
	return &Client{
		Conn:      conn,
		Send:      make(chan []byte, 256),
		Chatrooms: make(map[string]bool),
	}
}

// SendEvent sends an event to the client
func (c *Client) SendEvent(eventType EventType, payload interface{}) error {
	event := Event{
		Type:      eventType,
		Payload:   payload,
		Timestamp: TimeNow(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	c.Send <- data
	return nil
}

// Close closes the client's WebSocket connection and channels
func (c *Client) Close() {
	close(c.Send)
	c.Conn.Socket.Close()
}
