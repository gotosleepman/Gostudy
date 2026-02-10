package controllers

import (
	"fmt"
	"myblog-backend/database"
	"myblog-backend/handlers"
	"myblog-backend/middleware"
	"myblog-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController struct{}

// CreatePost 创建文章
func (pc *PostController) CreatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handlers.Unauthorized(c, "User not authenticated")
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.BadRequest(c, "Invalid request format")
		return
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	if err := database.DB.Create(&post).Error; err != nil {
		middleware.ErrorLogger(err, "Failed to create post")
		handlers.InternalServerError(c, "Failed to create post")
		return
	}

	// 获取完整的文章信息（包含用户）
	var createdPost models.Post
	if err := database.DB.Preload("User").First(&createdPost, post.ID).Error; err != nil {
		middleware.ErrorLogger(err, "Failed to load created post")
		handlers.InternalServerError(c, "Failed to retrieve created post")
		return
	}

	handlers.Created(c, "Post created successfully", createdPost)
}

// GetPosts 获取所有文章列表
func (pc *PostController) GetPosts(c *gin.Context) {
	var posts []models.Post

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// 查询文章，预加载用户信息
	query := database.DB.Model(&models.Post{}).Preload("User")

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取分页数据
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		middleware.ErrorLogger(err, "Failed to get posts")
		handlers.InternalServerError(c, "Failed to get posts")
		return
	}

	response := gin.H{
		"posts": posts,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	handlers.Success(c, "Posts retrieved successfully", response)
}

// GetPost 获取单个文章
func (pc *PostController) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handlers.BadRequest(c, "Invalid post ID")
		return
	}

	var post models.Post
	if err := database.DB.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handlers.NotFound(c, "Post not found")
		} else {
			middleware.ErrorLogger(err, "Failed to get post")
			handlers.InternalServerError(c, "Failed to get post")
		}
		return
	}

	handlers.Success(c, "Post retrieved successfully", post)
}

// UpdatePost 更新文章
func (pc *PostController) UpdatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handlers.Unauthorized(c, "User not authenticated")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handlers.BadRequest(c, "Invalid post ID")
		return
	}

	// 查找文章
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handlers.NotFound(c, "Post not found")
		} else {
			middleware.ErrorLogger(err, "Failed to find post for update")
			handlers.InternalServerError(c, "Failed to update post")
		}
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		handlers.Forbidden(c, "You are not the author of this post")
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.BadRequest(c, "Invalid request format")
		return
	}

	// 更新文章
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&post).Updates(updates).Error; err != nil {
			middleware.ErrorLogger(err, fmt.Sprintf("Failed to update post: %v", err))
			handlers.InternalServerError(c, "Failed to update post")
			return
		}
	}

	handlers.Success(c, "Post updated successfully", post)
}

// DeletePost 删除文章
func (pc *PostController) DeletePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handlers.Unauthorized(c, "User not authenticated")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handlers.BadRequest(c, "Invalid post ID")
		return
	}

	// 查找文章
	var post models.Post
	if err := database.DB.First(&post, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handlers.NotFound(c, "Post not found")
		} else {
			middleware.ErrorLogger(err, "Failed to find post for deletion")
			handlers.InternalServerError(c, "Failed to delete post")
		}
		return
	}

	// 检查是否是文章作者
	if post.UserID != userID.(uint) {
		handlers.Forbidden(c, "You are not the author of this post")
		return
	}

	// 删除文章（软删除）
	if err := database.DB.Delete(&post).Error; err != nil {
		middleware.ErrorLogger(err, fmt.Sprintf("Failed to delete post: %v", err))
		handlers.InternalServerError(c, "Failed to delete post")
		return
	}

	handlers.Success(c, "Post deleted successfully", nil)
}
