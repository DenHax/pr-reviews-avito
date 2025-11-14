package users

import storage "github.com/DenHax/pr-reviews-avito/internal/storage/postgres"

type UserStore struct {
	storage *storage.Storage
}

func NewUsersStorage(s *storage.Storage) *UserStore {
	return &UserStore{storage: s}
}

func (r *UserStore) GetAll() {
}
