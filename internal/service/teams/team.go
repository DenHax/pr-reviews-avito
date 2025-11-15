package teams

import (
	"errors"

	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	"github.com/DenHax/pr-reviews-avito/internal/repo"
	"github.com/DenHax/pr-reviews-avito/internal/service/users"
)

type TeamService struct {
	repo        repo.Teams
	userService users.UserService
}

func NewTeamService(repo repo.Teams, userService users.UserService) *TeamService {
	return &TeamService{repo: repo, userService: userService}
}

func (s *TeamService) CreateTeam(req models.TeamRequest) error {
	exists, err := s.repo.TeamExists(req.TeamName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("team already exists")
	}

	if err := s.repo.CreateTeam(req.TeamName); err != nil {
		return err
	}

	for _, member := range req.Members {
		user := &models.User{
			UserID:   member.UserID,
			Username: member.Username,
			TeamName: req.TeamName,
			IsActive: member.IsActive,
		}

		existingUser, err := s.userService.GetUser(member.UserID)
		if err != nil {
			return err
		}

		if existingUser == nil {
			if err := s.userService.CreateUser(user); err != nil {
				return err
			}
		} else {
			if err := s.userService.UpdateUser(user); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *TeamService) GetTeam(teamName string) (*models.TeamRequest, error) {
	exists, err := s.repo.TeamExists(teamName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("team not found")
	}

	users, err := s.repo.GetTeamMembers(teamName)
	if err != nil {
		return nil, err
	}

	members := make([]models.TeamMember, len(users))
	for i, user := range users {
		members[i] = models.TeamMember{
			UserID:   user.UserID,
			Username: user.Username,
			IsActive: user.IsActive,
		}
	}

	return &models.TeamRequest{
		TeamName: teamName,
		Members:  members,
	}, nil
}
