package users

import (
	"github.com/DenHax/pr-reviews-avito/internal/repo"
)

type UserService struct {
	repo repo.Users
}

func NewUserService(repo repo.Users) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAll() {
}
