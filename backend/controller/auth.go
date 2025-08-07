package controller

import (
	"backend/global"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	// 验证用户是否存在
	if err := global.DB.Where("email = ?", user.Email).First(&user).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "User already exists"})
		return
	}
	// 哈希密码
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password"})
		return
	}
	// 检查密码长度
	if len(user.Password) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Password must be at least 8 characters long"})
		return
	}
	user.Password = hashedPassword
	user.UserID = uuid.New().String()
	// 创建用户
	if err := global.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create user"})
		return
	}

	// JWT生成token
	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "User registered successfully", "token": token})
}

func Login(ctx *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	var user models.User
	if err := global.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}
	// 验证密码
	if !utils.VerifyPassword(user.Password, input.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid password"})
		return
	}
	// 生成token
	token, err := utils.GenerateToken(user.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate token"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successfully", "token": token})
}

func CheckToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token is required"})
		return
	}
	_, err := utils.VerifyToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid token"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Token is valid"})
}
