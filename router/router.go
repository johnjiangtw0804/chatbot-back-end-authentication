package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/johnjiangtw0804/chatbot-back-end-authentication/env"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/models"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/repository"
)

func RegisterRouter(env *config.Configuration, dbConnection *models.DBWrapper) (*gin.Engine, error) {
	router := gin.Default()
	router.Use(gin.Logger())   // log 每個請求
	router.Use(gin.Recovery()) // 保護程式不崩潰

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	v1 := router.Group("/api/v1")
	{
		userRepo := repository.NewUserRepository(dbConnection)
		registerUserRoutes(v1, userRepo)
	}

	return router, nil
}
