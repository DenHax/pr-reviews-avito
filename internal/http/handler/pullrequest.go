package handler

import (
	"log/slog"
	"net/http"

	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// CreatePR godoc
// @Summary Create a pull request
// @Description Create a new pull request with the specified details
// @Tags pull-requests
// @Accept json
// @Produce json
// @Param request body models.CreatePRRequest true "Pull request creation data"
// @Success 201 {object} models.PRResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /pullRequest/create [post]
func (h *Handler) CreatePR(c *gin.Context) {
	var req models.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid create PR request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
			},
		})
		return
	}

	pr, err := h.Services.CreatePR(req)
	if err != nil {
		slog.Error("Failed to create PR", "error", err, "pr_id", req.PullRequestID)

		switch err.Error() {
		case "author/team not found":
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "Author/team not found",
				},
			})
		case "PR already exists":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: ErrorDetail{
					Code:    "PR_EXISTS",
					Message: "PR id already exists",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: ErrorDetail{
					Code:    "INTERNAL_ERROR",
					Message: "Failed to create PR",
				},
			})
		}
		return
	}

	slog.Info("PR created successfully", "pr_id", req.PullRequestID, "author", req.AuthorID)
	c.JSON(http.StatusCreated, models.PRResponse{PR: pr})
}

// MergePR godoc
// @Summary Merge a pull request
// @Description Merge an existing pull request by its ID
// @Tags pull-requests
// @Accept json
// @Produce json
// @Param request body models.MergePRRequest true "Pull request merge data"
// @Success 200 {object} models.PRResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /pullRequest/merge [post]
func (h *Handler) MergePR(c *gin.Context) {
	var req models.MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid merge PR request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
			},
		})
		return
	}

	pr, err := h.Services.MergePR(req.PullRequestID)
	if err != nil {
		slog.Error("Failed to merge PR", "error", err, "pr_id", req.PullRequestID)

		if err.Error() == "PR not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "resource not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to merge PR",
			},
		})
		return
	}

	slog.Info("PR merged successfully", "pr_id", req.PullRequestID)
	c.JSON(http.StatusOK, models.PRResponse{PR: pr})
}

// Reassign godoc
// @Summary Reassign a reviewer
// @Description Reassign a reviewer from a pull request to another team member
// @Tags pull-requests
// @Accept json
// @Produce json
// @Param request body models.ReassignRequest true "Reassignment data"
// @Success 200 {object} models.ReassignResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /pullRequest/reassign [post]
func (h *Handler) Reassign(c *gin.Context) {
	var req models.ReassignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid reassign request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
			},
		})
		return
	}

	pr, replacedBy, err := h.Services.ReassignReviewer(req)
	if err != nil {
		slog.Error("Failed to reassign reviewer", "error", err, "pr_id", req.PullRequestID, "old_user", req.OldUserID)

		switch err.Error() {
		case "PR or user not found":
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "PR or user not found",
				},
			})
		case "cannot reassign on merged PR":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: ErrorDetail{
					Code:    "PR_MERGED",
					Message: "cannot reassign on merged PR",
				},
			})
		case "reviewer is not assigned to this PR":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_ASSIGNED",
					Message: "reviewer is not assigned to this PR",
				},
			})
		case "no active replacement candidate in team":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NO_CANDIDATE",
					Message: "no active replacement candidate in team",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: ErrorDetail{
					Code:    "INTERNAL_ERROR",
					Message: "Failed to reassign reviewer",
				},
			})
		}
		return
	}

	slog.Info("Reviewer reassigned successfully",
		"pr_id", req.PullRequestID,
		"old_user", req.OldUserID,
		"new_user", replacedBy,
	)
	c.JSON(http.StatusOK, models.ReassignResponse{
		PR:         pr,
		ReplacedBy: replacedBy,
	})
}
