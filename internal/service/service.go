package service

import (
	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	"github.com/DenHax/pr-reviews-avito/internal/repo"
	"github.com/DenHax/pr-reviews-avito/internal/service/pr"
	"github.com/DenHax/pr-reviews-avito/internal/service/teams"
	"github.com/DenHax/pr-reviews-avito/internal/service/users"
)

type Users interface {
	SetUserActive(req models.SetActiveRequest) (*models.User, error)
	GetUserReviews(userID string) (*models.ReviewListResponse, error)
}

type Teams interface {
	CreateTeam(req models.TeamRequest) error
	GetTeam(teamName string) (*models.TeamRequest, error)
}

type PullRequests interface {
	CreatePR(req models.CreatePRRequest) (*models.PullRequest, error)
	MergePR(pullRequestID string) (*models.PullRequest, error)
	ReassignReviewer(req models.ReassignRequest) (*models.PullRequest, string, error)
}

type Service struct {
	Users
	Teams
	PullRequests
}

func NewService(repos *repo.Repository) *Service {
	userService := users.NewUserService(repos.Users, repos.PullRequests)
	teamService := teams.NewTeamService(repos.Teams, *userService)
	prService := pr.NewPRService(repos.PullRequests, *userService)
	return &Service{
		Users:        userService,
		Teams:        teamService,
		PullRequests: prService,
	}
}
