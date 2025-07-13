package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	router.GET("/validate", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println(tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Token is valid ➜ return whatever info you want
		c.JSON(http.StatusOK, gin.H{
			"message": "Token is valid",
		})
	})

	v1 := router.Group("/api/v1")
	{
		userRepo := repository.NewUserRepository(dbConnection)
		registerUserRoutes(v1, userRepo, conf)
	}

	return router, nil
}
