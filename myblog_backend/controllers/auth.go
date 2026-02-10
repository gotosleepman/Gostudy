package controllers

import (
	"myblog-backend/database"
	"myblog-backend/handlers"
	"myblog-backend/middleware"
	"myblog-backend/models"
	"myblog-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct{}

// Register 用户注册
func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.BadRequest(c, "Invalid request format")
		return
	}

	// 验证密码强度
	if err := utils.ValidatePassword(req.Password); err != nil {
		handlers.BadRequest(c, err.Error())
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		handlers.BadRequest(c, "Username already exists")
		return
	}

	// 检查邮箱是否已存在
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		handlers.BadRequest(c, "Email already exists")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		middleware.ErrorLogger(err, "Failed to hash password")
		handlers.InternalServerError(c, "Failed to process password")
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		middleware.ErrorLogger(err, "Failed to create user")
		handlers.InternalServerError(c, "Failed to create user")
		return
	}

	handlers.Created(c, "User registered successfully", models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}

// Login 用户登录
func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.BadRequest(c, "Invalid request format")
		return
	}

	// 查找用户
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handlers.Unauthorized(c, "Invalid username or password")
		} else {
			middleware.ErrorLogger(err, "Failed to find user")
			handlers.InternalServerError(c, "Failed to process login")
		}
		return
	}

	// 验证密码
	if err := utils.CheckPassword(req.Password, user.Password); err != nil {
		handlers.Unauthorized(c, "Invalid username or password")
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		middleware.ErrorLogger(err, "Failed to generate token")
		handlers.InternalServerError(c, "Failed to generate authentication token")
		return
	}

	handlers.Success(c, "Login successful", models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	})
}

// GetProfile 获取当前用户信息
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handlers.Unauthorized(c, "User not authenticated")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handlers.NotFound(c, "User not found")
		} else {
			middleware.ErrorLogger(err, "Failed to get user profile")
			handlers.InternalServerError(c, "Failed to get user profile")
		}
		return
	}

	handlers.Success(c, "User profile retrieved", models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
