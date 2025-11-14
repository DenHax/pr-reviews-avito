package repo

import (
	"github.com/DenHax/pr-reviews-avito/internal/repo/users"
	storage "github.com/DenHax/pr-reviews-avito/internal/storage/postgres"
)

type Users interface {
	GetAll()
}

type Repository struct {
	Users
}

func NewRepository(s *storage.Storage) *Repository {
	return &Repository{
		Users: users.NewUsersStorage(s),
	}
}
