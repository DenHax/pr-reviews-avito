package teams

import (
	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	storage "github.com/DenHax/pr-reviews-avito/internal/storage/postgres"
)

type TeamStore struct {
	storage *storage.Storage
}

func NewTeamStorage(s *storage.Storage) *TeamStore {
	return &TeamStore{storage: s}
}

func (r *TeamStore) CreateTeam(teamName string) error {
	query := `INSERT INTO review_service.teams (team_name) VALUES ($1)`
	_, err := r.storage.DB.Exec(query, teamName)
	return err
}

func (r *TeamStore) TeamExists(teamName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM review_service.teams WHERE team_name = $1)`
	err := r.storage.DB.Get(&exists, query, teamName)
	return exists, err
}

func (r *TeamStore) GetTeamMembers(teamName string) ([]models.User, error) {
	var users []models.User
	query := `SELECT user_id, username, team_name, is_active, created_at, updated_at 
	          FROM review_service.users 
	          WHERE team_name = $1`
	err := r.storage.DB.Select(&users, query, teamName)
	return users, err
}
