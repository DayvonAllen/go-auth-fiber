package services

import (
	"example.com/app/domain"
	"example.com/app/repo"
	"fmt"
)

type AuthService interface {
	Login(username, password string) (*domain.User, string, error)
}

type DefaultAuthService struct {
	repo repo.AuthRepo
}

func (a DefaultAuthService) Login(username, password string) (*domain.User, string, error) {
	u, token, err := a.repo.Login(username, password)
	if err != nil {
		return nil, "", fmt.Errorf("error logging in: %w", err)
	}
	return u, token, nil
}

func NewAuthService(repository repo.AuthRepo) DefaultAuthService {
	return DefaultAuthService{repository}
}