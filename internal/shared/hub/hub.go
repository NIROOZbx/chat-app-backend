package hub

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Room struct {
	conns        map[int]*websocket.Conn
	mu           sync.Mutex
	ch           chan []byte
	ctx          context.Context
	cancel       context.CancelFunc
	cleanupTimer *time.Timer
}

type Manager struct {
	redis  *redis.Client
	rooms  map[int]*Room
	hubsMu sync.Mutex
}

func NewManager(redis *redis.Client) *Manager {
	return &Manager{
		redis: redis,
		rooms: make(map[int]*Room),
	}
}

func (m *Manager) GetOrCreate(roomID int) (*Room, bool) {
	m.hubsMu.Lock()
	defer m.hubsMu.Unlock()

	if room, exists := m.rooms[roomID]; exists {
		return room, false
	}

	ctx, cancel := context.WithCancel(context.Background())

	newRoom := &Room{
		conns:  make(map[int]*websocket.Conn),
		ch:     make(chan []byte, 1024),
		ctx:    ctx,
		cancel: cancel,
	}
	m.rooms[roomID] = newRoom
	return newRoom, true
}

func (r *Room) Add(userID int, conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.cleanupTimer != nil {
		r.cleanupTimer.Stop()
		r.cleanupTimer = nil
	}
	r.conns[userID] = conn
}

func (r *Room) Remove(userID int, roomID int, m *Manager) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.conns, userID)
	isEmpty := len(r.conns) == 0

	if isEmpty {
		r.cleanupTimer = time.AfterFunc(10*time.Second, func() {
			r.mu.Lock()
			currentCount := len(r.conns)
			r.mu.Unlock()

			if currentCount == 0 {
				m.Delete(roomID)
				log.Printf("[INFO] Room %d officially stopped and cleared from memory", roomID)
			}
		})
	}

}

func (m *Manager) Delete(roomID int) {
	m.hubsMu.Lock()
	defer m.hubsMu.Unlock()

	if rm, exists := m.rooms[roomID]; exists {
		rm.cancel()
		delete(m.rooms, roomID)

	}

}

func (r *Room) Snapshot() map[int]*websocket.Conn {

	r.mu.Lock()
	defer r.mu.Unlock()

	newMap := make(map[int]*websocket.Conn, len(r.conns))
	for id, conn := range r.conns {

		newMap[id] = conn
	}

	return newMap
}

func (r *Room) Channel() chan []byte     { return r.ch }
func (r *Room) Context() context.Context { return r.ctx }
func (r *Room) Done() <-chan struct{}    { return r.ctx.Done() }

func (r *Room) GetOnlineCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.conns)
}

func (m *Manager) GetOnlineCount(roomID int) int {
	m.hubsMu.Lock()
	room, exists := m.rooms[roomID]
	m.hubsMu.Unlock()

	if !exists {
		return 0
	}
	return room.GetOnlineCount()
}
