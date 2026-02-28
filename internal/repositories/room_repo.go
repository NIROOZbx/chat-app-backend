package repositories

import (
	"chat-app/internal/models"

	"github.com/jmoiron/sqlx"
)

type supabaseRepo struct {
	db *sqlx.DB
}

type RoomRepository interface {
	CreateRoom(room *models.Room, creatorID int) error
	GetAllRooms() ([]models.Room, error)
	GetJoinedRooms(userID int) ([]models.Room, error)
	GetRoomById(id int) (*models.Room, error)
	DeleteRoom(id int) error
	GetUserRole(roomID, userID int) (string, error)
}

func NewRoomRepo(db *sqlx.DB) RoomRepository {
	return &supabaseRepo{db: db}
}

func (s *supabaseRepo) CreateRoom(room *models.Room, creatorID int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO room (name, max_members, topic, description,image, is_private, invite_code)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;`

	err = tx.Get(room, query,
		room.Name,
		room.MaxMembers,
		room.Topic,
		room.Description,
		room.Image,
		room.IsPrivate,
		room.InviteCode,
	)
	if err != nil {
		return err
	}

	memberQuery := `INSERT INTO room_members (room_id, user_id, role) VALUES ($1, $2, $3)`
	if _, err := tx.Exec(memberQuery, room.ID, creatorID, "admin"); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *supabaseRepo) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	query := `SELECT * FROM room ORDER BY created_at DESC;`
	err := s.db.Select(&rooms, query)
	return rooms, err
}

func (s *supabaseRepo) GetJoinedRooms(userID int) ([]models.Room, error) {
	var rooms []models.Room
	query := `
		SELECT r.* 
		FROM room r 
		JOIN room_members rm ON r.id = rm.room_id 
		WHERE rm.user_id = $1 
		ORDER BY r.created_at DESC;`
	err := s.db.Select(&rooms, query, userID)
	return rooms, err
}

func (s *supabaseRepo) GetRoomById(id int) (*models.Room, error) {
	var room models.Room
	query := `SELECT * FROM room WHERE id = $1 LIMIT 1;`
	err := s.db.Get(&room, query, id)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *supabaseRepo) DeleteRoom(id int) error {
	query := `DELETE FROM room WHERE id = $1;`
	_, err := s.db.Exec(query, id)
	return err
}
