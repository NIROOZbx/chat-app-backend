package pubsub

import "time"

const (
	EventUserJoined = "room.user_joined"
	EventUserLeft   = "room.user_left"
	EventNewMessage = "message.new"
	EventTyping     = "room.typing"
	EventOnline     = "user.online"
	EventOffline    = "user.offline"
)

type BaseEvent struct {
	Type string `json:"type"`
}

type UserJoinedEvent struct {
	BaseEvent
	RoomID   int    `json:"RoomID"`
	UserID   int    `json:"UserID"`
	UserName string `json:"UserName"`
	Role     string `json:"Role"`
}

type UserLeaveEvent struct {
	BaseEvent
	RoomID   int    `json:"RoomID"`
	UserID   int    `json:"UserID"`
	UserName string `json:"UserName"`
}

type MessageEvent struct {
	BaseEvent
	RoomID    int       `json:"RoomID"`
	UserID    int       `json:"UserID"`
	UserName  string    `json:"UserName"`
	Content   string    `json:"Content"`
	MessageID int       `json:"MessageID"`
	SentAt    time.Time `json:"SentAt"`
}

type TypingEvent struct {
	BaseEvent
	RoomID   int    `json:"RoomID"`
	UserID   int    `json:"UserID"`
	UserName string `json:"UserName"`
}

type UserOfflineEvent struct {
	BaseEvent
	RoomID   int    `json:"RoomID"`
	UserID   int    `json:"UserID"`
	UserName string `json:"UserName"`
}
type UserOnlineEvent struct {
	BaseEvent
	RoomID   int    `json:"RoomID"`
	UserID   int    `json:"UserID"`
	UserName string `json:"UserName"`
}
