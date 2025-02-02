package routes

import (
	"github.com/gin-gonic/gin"
	"library-backend/internal/handlers"
	"library-backend/internal/middleware"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	authHandler := handlers.NewAuthHandler()
	bookHandler := handlers.NewBookHandler()

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.Logger())

	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	books := v1.Group("/books")
	{
		books.GET("", bookHandler.GetBooksBySubject)
		books.POST("/pickup", bookHandler.SubmitPickupSchedule)

		books.Use(middleware.ValidateTokenMiddleware())
		books.GET("/pickup-schedules", bookHandler.ListPickupSchedules)
	}

	return router
}
