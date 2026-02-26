package repositories

import (
	"chat-app/internal/models"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type JoinRoomRepo interface {
	JoinRoom(userID, roomID int, role string) error
	GetRoomByID(roomID int) (*models.Room, error)
	GetUserByID(userID int) (*models.User, error)
	IsRoomMember(userID, roomID int) (bool, error)
	LeaveRoom(roomID, userID int) error
	GetUserRole(roomID, userID int) (string, error)
}

func (s *supabaseRepo) GetUserRole(roomID, userID int) (string, error) {

	var role string
	query := `SELECT role from room_members where user_id=$1 and room_id=$2`

	err := s.db.Get(&role, query, userID, roomID)

	if err != nil {
		return "", err
	}

	return role, nil

}

func (s *supabaseRepo) JoinRoom(userID, roomID int, role string) error {
	query := `INSERT INTO room_members (room_id, user_id, role) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, roomID, userID, role)
	fmt.Println(err)
	return err
}

func (s *supabaseRepo) GetRoomByID(roomID int) (*models.Room, error) {
	var room models.Room
	query := `SELECT id, name, topic FROM room WHERE id = $1`
	err := s.db.Get(&room, query, roomID)
	return &room, err
}

func (s *supabaseRepo) IsRoomMember(userID, roomID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)`
	err := s.db.Get(&exists, query, roomID, userID)
	return exists, err
}

func (s *supabaseRepo) LeaveRoom(roomID, userID int) error {
	fmt.Println(roomID, userID)
	query := `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`
	result, err := s.db.Exec(query, roomID, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user is not member of this group")
	}
	return nil
}

func (s *supabaseRepo) GetUserByID(userID int) (*models.User, error) {
	var data models.User

	query := `SELECT * FROM users WHERE id = $1`

	err := s.db.Get(&data, query, userID)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func NewJoinRoomRepository(db *sqlx.DB) JoinRoomRepo {
	return &supabaseRepo{db: db}
}
