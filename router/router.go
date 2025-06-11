package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/johnjiangtw0804/chatbot-back-end-authentication/env"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/models"
)

func RegisterRouter(env *config.Configuration, dbConnection *models.DBWrapper) (*gin.Engine, error) {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	return r, nil
}
