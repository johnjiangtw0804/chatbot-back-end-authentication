package router

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/models"
	"github.com/johnjiangtw0804/chatbot-back-end-authentication/repository"
	"golang.org/x/crypto/bcrypt"
)

func registerUserRoutes(router *gin.RouterGroup, repo repository.UserRepository) {
	router.Group("/user")
	user := userEndPoint{repo: repo}
	router.POST("/register", user.registerHandler)
	router.PUT("/delete", user.deleteHandler)
}

type userEndPoint struct {
	repo repository.UserRepository
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (u *userEndPoint) registerHandler(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	existingUser, err := u.repo.FindByEmail(strings.ToLower(req.Email))
	if err == nil && existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := &models.User{
		Email:        strings.ToLower(req.Email),
		Name:         req.Name,
		PasswordHash: string(hashedPwd),
	}

	if err := u.repo.Create(newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 成功回傳
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    newUser.ID,
			"email": newUser.Email,
			"name":  newUser.Name,
		},
	})
}

type DeleteRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

func (u *userEndPoint) deleteHandler(ctx *gin.Context) {
	// 從 context 拿出目前登入的使用者 email（你可以存在 JWT claim 或 middleware 設定）
	email := ctx.GetString("userEmail") // 前提是你用 middleware 設定過
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
