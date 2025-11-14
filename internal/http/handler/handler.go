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

	health := router.Group("/health")
	{
		health.GET("", h.healthCheck)
		health.GET("/live", h.liveness)
		health.GET("/ready", h.readiness)
	}

	return router
}

func (h *Handler) redirectToSwagger(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
}
