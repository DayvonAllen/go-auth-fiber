package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUsers() (*[]domain.User, error)
	CreateUser(*domain.User) error
	UpdateUser(primitive.ObjectID, *domain.User) (*domain.User, error)
	GetUserByID(primitive.ObjectID) (*domain.User, error)
	DeleteByID(primitive.ObjectID) error
}

// the service has a dependency of the repo
type DefaultUserService struct {
	repo repo.UserRepo
}

func (s DefaultUserService) GetAllUsers() (*[]domain.User, error) {
	u, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}
	return  u, nil
}

func (s DefaultUserService) GetUserByID(id primitive.ObjectID) (*domain.User, error) {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (s DefaultUserService) UpdateUser(id primitive.ObjectID, user *domain.User) (*domain.User, error) {
	u, err := s.repo.UpdateByID(id, user)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}
	return u, nil
}

func (s DefaultUserService) DeleteByID(id primitive.ObjectID) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		return fmt.Errorf("error deleteing user: %w", err)
	}
	return nil
}

func (s DefaultUserService) CreateUser(user *domain.User) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	err := s.repo.Create(user)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func NewUserService(repository repo.UserRepo) DefaultUserService {
	return DefaultUserService{repository}
}