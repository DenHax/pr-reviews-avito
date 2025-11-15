package pr

import (
	"database/sql"
	"time"

	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	storage "github.com/DenHax/pr-reviews-avito/internal/storage/postgres"
)

type PRStore struct {
	storage *storage.Storage
}

func NewPRStorage(s *storage.Storage) *PRStore {
	return &PRStore{storage: s}
}

func (r *PRStore) CreatePR(pr *models.PullRequest) error {
	query := `INSERT INTO review_service.pull_requests 
	          (pull_request_id, pull_request_name, author_id, status, created_at) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.storage.DB.Exec(query, pr.PullRequestID, pr.PullRequestName, pr.AuthorID, pr.Status, pr.CreatedAt)
	return err
}

func (r *PRStore) GetPR(prID string) (*models.PullRequest, error) {
	var pr models.PullRequest
	query := `SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at 
	          FROM review_service.pull_requests 
	          WHERE pull_request_id = $1`
	err := r.storage.DB.Get(&pr, query, prID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &pr, err
}

func (r *PRStore) PRExists(prID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM review_service.pull_requests WHERE pull_request_id = $1)`
	err := r.storage.DB.Get(&exists, query, prID)
	return exists, err
}

func (r *PRStore) UpdatePRStatus(prID string, status string, mergedAt *time.Time) error {
	query := `UPDATE review_service.pull_requests 
	          SET status = $1, merged_at = $2 
	          WHERE pull_request_id = $3`
	_, err := r.storage.DB.Exec(query, status, mergedAt, prID)
	return err
}

func (r *PRStore) AssignReviewer(prID, userID string) error {
	query := `INSERT INTO review_service.pr_reviewers (pull_request_id, user_id) 
	          VALUES ($1, $2)`
	_, err := r.storage.DB.Exec(query, prID, userID)
	return err
}

func (r *PRStore) RemoveReviewer(prID, userID string) error {
	query := `DELETE FROM review_service.pr_reviewers 
	          WHERE pull_request_id = $1 AND user_id = $2`
	_, err := r.storage.DB.Exec(query, prID, userID)
	return err
}

func (r *PRStore) GetPRReviewers(prID string) ([]string, error) {
	var reviewers []string
	query := `SELECT user_id FROM review_service.pr_reviewers WHERE pull_request_id = $1`
	err := r.storage.DB.Select(&reviewers, query, prID)
	return reviewers, err
}

func (r *PRStore) GetPRsByReviewer(userID string) ([]models.PullRequest, error) {
	var prs []models.PullRequest
	query := `SELECT p.pull_request_id, p.pull_request_name, p.author_id, p.status, p.created_at, p.merged_at
	          FROM review_service.pull_requests p
	          JOIN review_service.pr_reviewers pr ON p.pull_request_id = pr.pull_request_id
	          WHERE pr.user_id = $1`
	err := r.storage.DB.Select(&prs, query, userID)
	return prs, err
}

func (r *PRStore) GetReviewerCount(prID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM review_service.pr_reviewers WHERE pull_request_id = $1`
	err := r.storage.DB.Get(&count, query, prID)
	return count, err
}

func (r *PRStore) IsUserReviewer(prID, userID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM review_service.pr_reviewers WHERE pull_request_id = $1 AND user_id = $2)`
	err := r.storage.DB.Get(&exists, query, prID, userID)
	return exists, err
}
