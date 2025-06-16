package router

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/env"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/models"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/repository"
)

func RegisterRouter(conf *env.Configuration, dbConnection *models.DBWrapper) (*gin.Engine, error) {
	router := gin.Default()
	router.Use(gin.Logger())   // log 每個請求
	router.Use(gin.Recovery()) // 保護程式不崩潰

	router.Use(cors.New(cors.Config{AllowOrigins: []string{conf.AppFrontendURL},
		MaxAge:           86400,
		AllowMethods:     []string{"POST, GET, OPTIONS, PUT, DELETE, UPDATE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true, // to allow browsers 自帶 credentials
		ExposeHeaders:    []string{"Content-Length"},
		// cache control header, some static assets can be cached in the browser
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	v1 := router.Group("/api/v1")
	{
		userRepo := repository.NewUserRepository(dbConnection)
		registerUserRoutes(v1, userRepo, conf)
	}

	return router, nil
}
