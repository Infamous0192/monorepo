package entity

// ChatroomType represents the type of chatroom
type ChatroomType string

const (
	ChatroomTypePublic  ChatroomType = "public"
	ChatroomTypePrivate ChatroomType = "private"
	ChatroomTypeSquad   ChatroomType = "squad"
)

// ParticipantRole represents the role of a participant in a chatroom
type ParticipantRole string

const (
	ParticipantRoleMember     ParticipantRole = "member"
	ParticipantRoleAdmin      ParticipantRole = "admin"
	ParticipantRoleSuperAdmin ParticipantRole = "super_admin"
)

// Chatroom represents a chatroom for group or direct conversations
type Chatroom struct {
	ID                   string                `bson:"_id,omitempty" json:"id,omitempty"`
	Name                 string                `bson:"name" json:"name"`
	IsGroup              bool                  `bson:"isGroup" json:"isGroup"`
	Type                 ChatroomType          `bson:"type" json:"type"`
	CreatedTimestamp     int64                 `bson:"createdTimestamp" json:"createdTimestamp"`
	LastMessage          *string               `bson:"lastMessage" json:"lastMessage"`
	LastSender           *string               `bson:"lastSender" json:"lastSender"`
	LastMessageTimestamp *int64                `bson:"lastMessageTimestamp" json:"lastMessageTimestamp"`
	MessagesCount        int                   `bson:"messagesCount" json:"messagesCount"`
	Participants         []ChatroomParticipant `bson:"participants" json:"participants"`
}

// ChatroomParticipant represents a participant in a chatroom
type ChatroomParticipant struct {
	ID                  string          `bson:"_id,omitempty" json:"id,omitempty"`
	User                string          `bson:"user" json:"user"`
	Role                ParticipantRole `bson:"role" json:"role"`
	JoinedTimestamp     int64           `bson:"joinedTimestamp" json:"joinedTimestamp"`
	MutedUntilTimestamp *int64          `bson:"mutedUntilTimestamp,omitempty" json:"mutedUntilTimestamp,omitempty"`
}

// ChatroomParticipantPopulated represents a participant in a chatroom with populated user
type ChatroomParticipantPopulated struct {
	ID                  string          `bson:"_id,omitempty" json:"id,omitempty"`
	User                User            `bson:"user" json:"user"` // Populated User object
	Role                ParticipantRole `bson:"role" json:"role"`
	JoinedTimestamp     int64           `bson:"joinedTimestamp" json:"joinedTimestamp"`
	MutedUntilTimestamp *int64          `bson:"mutedUntilTimestamp,omitempty" json:"mutedUntilTimestamp,omitempty"`
}

// ChatroomPopulated represents a chatroom with populated user references
type ChatroomPopulated struct {
	ID                   string                         `bson:"_id,omitempty" json:"id,omitempty"`
	Name                 string                         `bson:"name" json:"name"`
	IsGroup              bool                           `bson:"isGroup" json:"isGroup"`
	Type                 ChatroomType                   `bson:"type" json:"type"`
	CreatedTimestamp     int64                          `bson:"createdTimestamp" json:"createdTimestamp"`
	LastMessage          *string                        `bson:"lastMessage" json:"lastMessage"`
	LastSender           *User                          `bson:"lastSender" json:"lastSender"` // Populated User object
	LastMessageTimestamp *int64                         `bson:"lastMessageTimestamp" json:"lastMessageTimestamp"`
	MessagesCount        int                            `bson:"messagesCount" json:"messagesCount"`
	Participants         []ChatroomParticipantPopulated `bson:"participants" json:"participants"` // Array of populated participants
}
