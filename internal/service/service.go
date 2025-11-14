package service

import (
	"github.com/DenHax/pr-reviews-avito/internal/repo"
	"github.com/DenHax/pr-reviews-avito/internal/service/users"
)

type Users interface {
	GetAll()
}

type Service struct {
	Users
}

func NewService(repos *repo.Repository) *Service {
	return &Service{
		Users: users.NewUserService(repos.Users),
	}
}
