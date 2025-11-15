package users

import (
	"errors"

	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	"github.com/DenHax/pr-reviews-avito/internal/repo"
)

type UserService struct {
	repo   repo.Users
	prRepo repo.PullRequests
}

func NewUserService(repo repo.Users, prRepo repo.PullRequests) *UserService {
	return &UserService{repo: repo, prRepo: prRepo}
}

func (s *UserService) SetUserActive(req models.SetActiveRequest) (*models.User, error) {
	user, err := s.repo.GetUser(req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := s.repo.SetUserActive(req.UserID, req.IsActive); err != nil {
		return nil, err
	}

	user.IsActive = req.IsActive
	return user, nil
}

func (s *UserService) GetUserReviews(userID string) (*models.ReviewListResponse, error) {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	prs, err := s.prRepo.GetPRsByReviewer(userID)
	if err != nil {
		return nil, err
	}

	prShorts := make([]models.PullRequestShort, len(prs))
	for i, pr := range prs {
		prShorts[i] = models.PullRequestShort{
			PullRequestID:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorID:        pr.AuthorID,
			Status:          pr.Status,
		}
	}

	return &models.ReviewListResponse{
		UserID:       userID,
		PullRequests: prShorts,
	}, nil
}

func (s *UserService) GetUser(userID string) (*models.User, error) {
	return s.repo.GetUser(userID)
}

func (s *UserService) GetActiveTeamMembers(teamName string, excludeUserID string) ([]models.User, error) {
	return s.repo.GetActiveTeamMembers(teamName, excludeUserID)
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.repo.CreateUser(user)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.repo.UpdateUser(user)
}
