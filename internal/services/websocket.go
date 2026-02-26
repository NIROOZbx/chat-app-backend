package services

import (
	"chat-app/internal/shared/hub"
	"chat-app/internal/shared/logger"
	"chat-app/internal/shared/pubsub"
	"context"

	"github.com/gorilla/websocket"
)

type wsService struct {
	pubsub  *pubsub.PubSub
	manager *hub.Manager
	log     *logger.Logger
}

type WSService interface {
	Connect(userName string, roomID, userID int, conn *websocket.Conn)
	Disconnect(username string, roomID, userID int)
	UserTyping(ctx context.Context, userName string, userID, roomID int)
}

func (ws *wsService) Connect(userName string, roomID, userID int, conn *websocket.Conn) {

	room, created := ws.manager.GetOrCreate(roomID)

	if created {
		ws.pubsub.Subscribe(room.Context(), roomID, room.Channel())
		go ws.broadcast(room, roomID)

	}
	room.Add(userID, conn)
	payload := pubsub.UserOnlineEvent{
		BaseEvent: pubsub.BaseEvent{Type: pubsub.EventOnline},
		RoomID:    roomID,
		UserID:    userID,
		UserName:  userName,
	}
	ws.pubsub.Publish(room.Context(), roomID, payload)

}

func (ws *wsService) Disconnect(username string, roomID, userID int) {

	room, created := ws.manager.GetOrCreate(roomID)

	if created {
		ws.manager.Delete(roomID)
		return
	}
	payload := pubsub.UserOfflineEvent{
		BaseEvent: pubsub.BaseEvent{Type: pubsub.EventOffline},
		RoomID:    roomID,
		UserID:    userID,
		UserName:  username,
	}

	ws.pubsub.Publish(context.Background(), roomID, payload)

	room.Remove(userID, roomID, ws.manager)

}

func (ws *wsService) UserTyping(ctx context.Context, userName string, userID, roomID int) {

	payload := pubsub.TypingEvent{
		BaseEvent: pubsub.BaseEvent{Type: pubsub.EventTyping},
		RoomID:    roomID,
		UserID:    userID,
		UserName:  userName,
	}
	ws.pubsub.Publish(ctx, roomID, payload)
}
func (ws *wsService) broadcast(room *hub.Room, roomID int) {

	for {
		select {
		case msg, ok := <-room.Channel():
			if !ok {
				return
			}
			for userID, conn := range room.Snapshot() {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					ws.log.Error("write error user %d: %v", userID, err)
					conn.Close()
					room.Remove(userID, roomID, ws.manager)
				}
			}
		case <-room.Done():
			return

		}
	}

}

func NewWsService(pubsub *pubsub.PubSub, manager *hub.Manager, log *logger.Logger) WSService {
	return &wsService{
		pubsub:  pubsub,
		manager: manager,
		log:     log,
	}
}
