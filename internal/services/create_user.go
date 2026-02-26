package services

import (
	"chat-app/internal/models"
	"chat-app/internal/repositories"
	"chat-app/internal/shared/request"
)

type CreateService interface {
	CreateUser(req request.CreateUser, img string) (*models.User, error)
	CheckUser(name string) (*models.User, error)
	GetMe(id int) (*models.User, error)
}

type userService struct {
	Repo repositories.UserRepo
}

func (s *userService) CreateUser(req request.CreateUser, img string) (*models.User, error) {
	userModel := &models.User{
		UserName:     req.UserName,
		ProfileImage: img,
	}

	return s.Repo.CreateUser(userModel)
}

func (s *userService) CheckUser(name string) (*models.User, error) {
	return s.Repo.GetUserByName(name)
}

func (s *userService) GetMe(id int) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}

func NewUserService(repo repositories.UserRepo) CreateService {
	return &userService{Repo: repo}
}
