DROP TRIGGER IF EXISTS trigger_max_reviewers ON review_service.pr_reviewers;
DROP TRIGGER IF EXISTS update_users_updated_at ON review_service.users;
DROP FUNCTION IF EXISTS review_service.check_max_reviewers;
DROP FUNCTION IF EXISTS review_service.update_updated_at_column;
