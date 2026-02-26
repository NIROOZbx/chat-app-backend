package repositories

import (
	"chat-app/internal/models"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	CreateUser(m *models.User) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
}

func NewUserRepository(db *sqlx.DB) UserRepo {
	return &supabaseRepo{db: db}
}

func (s *supabaseRepo) CreateUser(m *models.User) (*models.User, error) {
	var data models.User
	query := `INSERT INTO users (user_name, profile_image) VALUES ($1, $2) RETURNING *;`
	err := s.db.Get(&data, query, m.UserName, m.ProfileImage)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *supabaseRepo) GetUserByName(name string) (*models.User, error) {

	var data models.User

	query := `select * from users where user_name=$1`

	err := s.db.Get(&data, query, name)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("user was not found")
	}
	return &data, nil

}
