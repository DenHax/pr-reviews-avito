package handler

import (
	"net/http"
	"time"

	"github.com/DenHax/pr-reviews-avito/internal/middleware"
	"github.com/DenHax/pr-reviews-avito/internal/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var startTime = time.Now()

type Handler struct {
	Services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Services: services,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.New()

	router.Use(middleware.CORSMIddleware())
	router.GET("/swagger", h.redirectToSwagger)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", h.CheckHealth)

	teamGroup := router.Group("/team")
	{
		teamGroup.POST("/add", h.AddTeam)
		teamGroup.GET("/get", h.GetTeam)
	}

	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/setIsActive", h.SetIsActive)
		usersGroup.GET("/getReview", h.GetReview)
	}

	pullRequestGroup := router.Group("/pullRequest")
	{
		pullRequestGroup.POST("/create", h.CreatePR)
		pullRequestGroup.POST("/merge", h.MergePR)
		pullRequestGroup.POST("/reassign", h.Reassign)
	}

	return router
}

func (h *Handler) redirectToSwagger(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}
