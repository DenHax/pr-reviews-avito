package handler

import (
	"log/slog"
	"net/http"

	"github.com/DenHax/pr-reviews-avito/internal/domain/models"
	"github.com/gin-gonic/gin"
)

// AddTeam godoc
// @Summary Create a new team
// @Description Create a new team with the specified name and members
// @Tags teams
// @Accept json
// @Produce json
// @Param request body models.TeamRequest true "Team creation data"
// @Success 201 {object} models.TeamResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /team/add [post]
func (h *Handler) AddTeam(c *gin.Context) {
	var req models.TeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid team request", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
			},
		})
		return
	}

	if err := h.Services.CreateTeam(req); err != nil {
		slog.Error("Failed to create team", "error", err, "team_name", req.TeamName)

		if err.Error() == "team already exists" {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: ErrorDetail{
					Code:    "TEAM_EXISTS",
					Message: "team_name already exists",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to create team",
			},
		})
		return
	}

	slog.Info("Team created successfully", "team_name", req.TeamName)
	c.JSON(http.StatusCreated, models.TeamResponse{Team: req})
}

// GetTeam godoc
// @Summary Get team information
// @Description Retrieve team details by team name
// @Tags teams
// @Accept json
// @Produce json
// @Param team_name query string true "Team name"
// @Success 200 {object} models.TeamResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /team/get [get]
func (h *Handler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		slog.Warn("Missing team_name query parameter")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: ErrorDetail{
				Code:    "MISSING_PARAMETER",
				Message: "team_name query parameter is required",
			},
		})
		return
	}

	team, err := h.Services.GetTeam(teamName)
	if err != nil {
		slog.Warn("Team not found", "team_name", teamName, "error", err)
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: ErrorDetail{
				Code:    "NOT_FOUND",
				Message: "resource not found",
			},
		})
		return
	}

	slog.Debug("Team retrieved", "team_name", teamName)
	c.JSON(http.StatusOK, team)
}
