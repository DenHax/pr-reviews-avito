package pr

import (
	"errors"
	"math/rand"
	"time"

	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	"github.com/DenHax/pr-reviews-avito/internal/repo"
	"github.com/DenHax/pr-reviews-avito/internal/service/users"
)

type PRService struct {
	repo        repo.PullRequests
	userService users.UserService
}

func NewPRService(repo repo.PullRequests, userService users.UserService) *PRService {
	return &PRService{repo: repo, userService: userService}
}

func (s *PRService) CreatePR(req models.CreatePRRequest) (*models.PullRequest, error) {
	exists, err := s.repo.PRExists(req.PullRequestID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("PR already exists")
	}

	author, err := s.userService.GetUser(req.AuthorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("author/team not found")
	}

	potentialReviewers, err := s.userService.GetActiveTeamMembers(author.TeamName, req.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewers := s.selectRandomReviewers(potentialReviewers, 2)

	now := time.Now()
	pr := &models.PullRequest{
		PullRequestID:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorID,
		Status:          "OPEN",
		CreatedAt:       now,
		MergedAt:        nil,
	}

	if err := s.repo.CreatePR(pr); err != nil {
		return nil, err
	}

	for _, reviewer := range reviewers {
		if err := s.repo.AssignReviewer(req.PullRequestID, reviewer.UserID); err != nil {
			return nil, err
		}
	}

	return pr, nil
}

func (s *PRService) MergePR(pullRequestID string) (*models.PullRequest, error) {
	pr, err := s.repo.GetPR(pullRequestID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, errors.New("PR not found")
	}

	if pr.Status == "MERGED" {
		return pr, nil
	}

	now := time.Now()
	if err := s.repo.UpdatePRStatus(pullRequestID, "MERGED", &now); err != nil {
		return nil, err
	}

	pr.Status = "MERGED"
	pr.MergedAt = &now

	return pr, nil
}

func (s *PRService) ReassignReviewer(req models.ReassignRequest) (*models.PullRequest, string, error) {
	pr, err := s.repo.GetPR(req.PullRequestID)
	if err != nil {
		return nil, "", err
	}
	if pr == nil {
		return nil, "", errors.New("PR or user not found")
	}

	if pr.Status == "MERGED" {
		return nil, "", errors.New("cannot reassign on merged PR")
	}

	isReviewer, err := s.repo.IsUserReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		return nil, "", err
	}
	if !isReviewer {
		return nil, "", errors.New("reviewer is not assigned to this PR")
	}

	oldUser, err := s.userService.GetUser(req.OldUserID)
	if err != nil {
		return nil, "", err
	}
	if oldUser == nil {
		return nil, "", errors.New("PR or user not found")
	}

	candidates, err := s.userService.GetActiveTeamMembers(oldUser.TeamName, req.OldUserID)
	if err != nil {
		return nil, "", err
	}

	currentReviewers, err := s.repo.GetPRReviewers(req.PullRequestID)
	if err != nil {
		return nil, "", err
	}

	availableCandidates := s.filterExistingReviewers(candidates, currentReviewers)

	if len(availableCandidates) == 0 {
		return nil, "", errors.New("no active replacement candidate in team")
	}

	newReviewer := availableCandidates[rand.Intn(len(availableCandidates))]

	if err := s.repo.RemoveReviewer(req.PullRequestID, req.OldUserID); err != nil {
		return nil, "", err
	}

	if err := s.repo.AssignReviewer(req.PullRequestID, newReviewer.UserID); err != nil {
		return nil, "", err
	}

	return pr, newReviewer.UserID, nil
}

func (s *PRService) selectRandomReviewers(users []models.User, max int) []models.User {
	if len(users) == 0 {
		return []models.User{}
	}

	shuffled := make([]models.User, len(users))
	copy(shuffled, users)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	if len(shuffled) > max {
		return shuffled[:max]
	}
	return shuffled
}

func (s *PRService) filterExistingReviewers(candidates []models.User, existingReviewers []string) []models.User {
	reviewerSet := make(map[string]bool)
	for _, reviewer := range existingReviewers {
		reviewerSet[reviewer] = true
	}

	filtered := make([]models.User, 0)
	for _, candidate := range candidates {
		if !reviewerSet[candidate.UserID] {
			filtered = append(filtered, candidate)
		}
	}
	return filtered
}
