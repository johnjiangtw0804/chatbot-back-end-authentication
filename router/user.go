package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/env"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/models"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/repository"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/router/middleware"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/utils"
	"golang.org/x/crypto/bcrypt"
)

func registerUserRoutes(router *gin.RouterGroup, repo repository.UserRepository, config *env.Configuration) {
	// public API
	userRouter := router.Group("/user")
	useEndPoint := userEndPoint{repo: repo, config: config}
	router.POST("/register", useEndPoint.registerHandler)

	// authorized API
	userRouter.Use(middleware.JWTMiddleware([]byte(config.JWTSecret)))
	userRouter.DELETE("/delete", useEndPoint.deleteHandler)
}

type userEndPoint struct {
	repo   repository.UserRepository
	config *env.Configuration
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *userEndPoint) registerHandler(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if len(req.Password) < 6 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: Password Length < 6"})
		return
	}

	existingUser, err := u.repo.FindByEmail(strings.ToLower(req.Email))
	if err == nil && existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Default().Println("Failed to hash password")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := &models.User{
		Email:        strings.ToLower(req.Email),
		Name:         req.Name,
		PasswordHash: string(hashedPwd),
	}

	if err := u.repo.Create(newUser); err != nil {
		log.Default().Println("Failed to create user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// generate a token and return to the user
	log.Default().Println(fmt.Sprintf("%s %d %s", u.config.JWTSecret, newUser.ID, newUser.Email))
	token, err := utils.GenerateJWT(u.config.JWTSecret, newUser.ID, newUser.Email)
	if err != nil {
		log.Default().Println("Failed to generate JWT " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
	})
}

type DeleteRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func (u *userEndPoint) deleteHandler(ctx *gin.Context) {

	// 從 context 拿出目前登入的使用者 email（你可以存在 JWT claim 或 middleware 設定）
	email := ctx.GetString("email") // 前提是你用 middleware 設定過
	if email == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 先查出 user
	user, err := u.repo.FindByEmail(email)
	if err != nil || user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 執行刪除
	err = u.repo.Delete(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed: " + err.Error()})
		return
	}

	log.Printf("User deleted: %s (%s)\n", user.Name, user.Email)
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
