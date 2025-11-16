package handler

import (
	"log/slog"
	"net/http"

	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// SetIsActive godoc
// @Summary Update user active status
// @Description Set or update the active status of a user
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.SetActiveRequest true "User active status data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/setIsActive [post]
func (h *Handler) SetIsActive(c *gin.Context) {
	var req models.SetActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid set active request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
			},
		})
		return
	}

	user, err := h.Services.SetUserActive(req)
	if err != nil {
		slog.Error("Failed to set user active", "error", err, "user_id", req.UserID)

		if err.Error() == "user not found" {
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
				Message: "Failed to update user",
			},
		})
		return
	}

	if user == nil {
		slog.Error("User is nil after update", "user_id", req.UserID)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to update user",
			},
		})
		return
	}

	slog.Info("User active status updated", "user_id", req.UserID, "is_active", req.IsActive)
	c.JSON(http.StatusOK, models.UserResponse{User: *user})
}

// GetReview godoc
// @Summary Get user reviews
// @Description Retrieve all pull request reviews assigned to a user
// @Tags users
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} models.ReviewListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/getReview [get]
func (h *Handler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		slog.Warn("Missing user_id query parameter")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "MISSING_PARAMETER",
				Message: "user_id query parameter is required",
			},
		})
		return
	}

	reviews, err := h.Services.GetUserReviews(userID)
	if err != nil {
		slog.Error("Failed to get user reviews", "error", err, "user_id", userID)

		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: ErrorDetail{
					Code:    "NOT_FOUND",
					Message: "User not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to retrieve reviews",
			},
		})
		return
	}

	slog.Debug("User reviews retrieved", "user_id", userID, "pr_count", len(reviews.PullRequests))
	c.JSON(http.StatusOK, reviews)
}
