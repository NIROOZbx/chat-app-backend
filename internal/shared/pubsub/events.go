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
	RoomID   int    `json:"room_id"`
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Role     string `json:"role"`
}

type UserLeaveEvent struct {
	BaseEvent
	RoomID   int    `json:"room_id"`
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

type MessageEvent struct {
	BaseEvent
	RoomID    int       `json:"room_id"`
	UserID    int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	MessageID int       `json:"message_id"`
	SentAt    time.Time `json:"sent_at"`
}

type TypingEvent struct {
	BaseEvent
	RoomID   int    `json:"room_id"`
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

type UserOfflineEvent struct {
	BaseEvent
	RoomID int `json:"room_id"`
	UserID int `json:"user_id"`
	UserName string `json:"user_name"`
}
type UserOnlineEvent struct {
	BaseEvent
	RoomID int `json:"room_id"`
	UserID int `json:"user_id"`
	UserName string `json:"user_name"`

}
