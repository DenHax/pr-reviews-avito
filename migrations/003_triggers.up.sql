CREATE OR REPLACE FUNCTION review_service.check_max_reviewers()
RETURNS TRIGGER AS $$
BEGIN
    IF (
        SELECT COUNT(*) 
        FROM review_service.pr_reviewers 
        WHERE pull_request_id = NEW.pull_request_id
    ) > 2 THEN
        RAISE EXCEPTION 'Cannot assign more than 2 reviewers to a PR';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_max_reviewers
    BEFORE INSERT ON review_service.pr_reviewers
    FOR EACH ROW
    EXECUTE FUNCTION review_service.check_max_reviewers();

CREATE OR REPLACE FUNCTION review_service.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON review_service.users
    FOR EACH ROW
    EXECUTE FUNCTION review_service.update_updated_at_column();
