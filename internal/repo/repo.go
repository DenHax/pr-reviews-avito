package repo

import (
	"time"

	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	"github.com/DenHax/pr-reviews-avito/internal/repo/pr"
	"github.com/DenHax/pr-reviews-avito/internal/repo/teams"
	"github.com/DenHax/pr-reviews-avito/internal/repo/users"
	storage "github.com/DenHax/pr-reviews-avito/internal/storage/postgres"
)

type Users interface {
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetUser(userID string) (*models.User, error)
	GetActiveTeamMembers(teamName string, excludeUserID string) ([]models.User, error)
	SetUserActive(userID string, isActive bool) error
}

type Teams interface {
	CreateTeam(teamName string) error
	TeamExists(teamName string) (bool, error)
	GetTeamMembers(teamName string) ([]models.User, error)
}

type PullRequests interface {
	CreatePR(pr *models.PullRequest) error
	GetPR(prID string) (*models.PullRequest, error)
	PRExists(prID string) (bool, error)
	UpdatePRStatus(prID string, status string, mergedAt *time.Time) error
	AssignReviewer(prID, userID string) error
	RemoveReviewer(prID, userID string) error
	GetPRReviewers(prID string) ([]string, error)
	GetPRsByReviewer(userID string) ([]models.PullRequest, error)
	GetReviewerCount(prID string) (int, error)
	IsUserReviewer(prID, userID string) (bool, error)
}

type Repository struct {
	Users
	Teams
	PullRequests
}

func NewRepository(s *storage.Storage) *Repository {
	return &Repository{
		Users:        users.NewUsersStorage(s),
		Teams:        teams.NewTeamStorage(s),
		PullRequests: pr.NewPRStorage(s),
	}
}
