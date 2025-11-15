package users

import (
	"database/sql"

	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	storage "github.com/DenHax/pr-reviews-avito/internal/storage/postgres"
)

type UserStore struct {
	storage *storage.Storage
}

func NewUsersStorage(s *storage.Storage) *UserStore {
	return &UserStore{storage: s}
}

func (r *UserStore) CreateUser(user *models.User) error {
	query := `INSERT INTO review_service.users (user_id, username, team_name, is_active) 
	          VALUES ($1, $2, $3, $4)`
	_, err := r.storage.DB.Exec(query, user.UserID, user.Username, user.TeamName, user.IsActive)
	return err
}

func (r *UserStore) UpdateUser(user *models.User) error {
	query := `UPDATE review_service.users 
	          SET username = $1, team_name = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP 
	          WHERE user_id = $4`
	_, err := r.storage.DB.Exec(query, user.Username, user.TeamName, user.IsActive, user.UserID)
	return err
}

func (r *UserStore) GetUser(userID string) (*models.User, error) {
	var user models.User
	query := `SELECT user_id, username, team_name, is_active, created_at, updated_at 
	          FROM review_service.users 
	          WHERE user_id = $1`
	err := r.storage.DB.Get(&user, query, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserStore) GetActiveTeamMembers(teamName string, excludeUserID string) ([]models.User, error) {
	var users []models.User
	query := `SELECT user_id, username, team_name, is_active, created_at, updated_at 
	          FROM review_service.users 
	          WHERE team_name = $1 AND is_active = true AND user_id != $2`
	err := r.storage.DB.Select(&users, query, teamName, excludeUserID)
	return users, err
}

func (r *UserStore) SetUserActive(userID string, isActive bool) error {
	query := `UPDATE review_service.users 
	          SET is_active = $1, updated_at = CURRENT_TIMESTAMP 
	          WHERE user_id = $2`
	_, err := r.storage.DB.Exec(query, isActive, userID)
	return err
}
