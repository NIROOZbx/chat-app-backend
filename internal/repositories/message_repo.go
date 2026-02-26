package repositories

import (
	"chat-app/internal/models"


	"github.com/jmoiron/sqlx"
)

type MessageRepo interface {
	Save(m *models.Message) (*models.Message, error)
	GetByRoomID(roomID int, limit int, offset int) ([]models.Message, error)
}

type messageRepo struct {
	db *sqlx.DB
}

func NewMessageRepo(db *sqlx.DB) MessageRepo {
	return &messageRepo{db: db}
}

func (r *messageRepo) Save(m *models.Message) (*models.Message, error) {
	query := `
        INSERT INTO messages (room_id, user_id, content)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
	if err := r.db.QueryRowx(query, m.RoomID, m.UserID, m.Content).
		Scan(&m.ID, &m.CreatedAt); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *messageRepo) GetByRoomID(roomID, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	query := `
        SELECT 
        m.id, m.room_id, m.user_id, m.content, m.created_at, u.user_name
        FROM messages AS m
        JOIN users AS u ON m.user_id = u.id
        WHERE m.room_id = $1
        ORDER BY m.created_at DESC
        LIMIT $2 OFFSET $3
    `

    
    if err := r.db.Select(&messages, query, roomID, limit, offset); err != nil {
        return nil, err
    }

    
    return messages, nil
}
