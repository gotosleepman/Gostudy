package controllers

import (
	"myblog-backend/database"
	"myblog-backend/handlers"
	"myblog-backend/middleware"
	"myblog-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CommentController struct{}

// CreateComment 创建评论
func (cc *CommentController) CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handlers.Unauthorized(c, "User not authenticated")
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlers.BadRequest(c, "Invalid request format")
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := database.DB.First(&post, req.PostID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handlers.NotFound(c, "Post not found")
		} else {
			middleware.ErrorLogger(err, "Failed to find post for comment")
			handlers.InternalServerError(c, "Failed to create comment")
		}
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID.(uint),
		PostID:  req.PostID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		middleware.ErrorLogger(err, "Failed to create comment")
		handlers.InternalServerError(c, "Failed to create comment")
		return
	}

	// 获取完整的评论信息（包含用户）
	var createdComment models.Comment
	if err := database.DB.Preload("User").First(&createdComment, comment.ID).Error; err != nil {
		middleware.ErrorLogger(err, "Failed to load created comment")
		handlers.InternalServerError(c, "Failed to retrieve created comment")
		return
	}

	handlers.Created(c, "Comment created successfully", createdComment)
}

// GetComments 获取文章的所有评论
func (cc *CommentController) GetComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		handlers.BadRequest(c, "Invalid post ID")
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handlers.NotFound(c, "Post not found")
		} else {
			middleware.ErrorLogger(err, "Failed to find post for comments")
			handlers.InternalServerError(c, "Failed to get comments")
		}
		return
	}

	var comments []models.Comment
	if err := database.DB.Preload("User").Where("post_id = ?", postID).Order("created_at DESC").Find(&comments).Error; err != nil {
		middleware.ErrorLogger(err, "Failed to get comments")
		handlers.InternalServerError(c, "Failed to get comments")
		return
	}

	handlers.Success(c, "Comments retrieved successfully", comments)
}

// DeleteComment 删除评论
func (cc *CommentController) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handlers.Unauthorized(c, "User not authenticated")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handlers.BadRequest(c, "Invalid comment ID")
		return
	}

	// 查找评论
	var comment models.Comment
	if err := database.DB.First(&comment, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handlers.NotFound(c, "Comment not found")
		} else {
			middleware.ErrorLogger(err, "Failed to find comment for deletion")
			handlers.InternalServerError(c, "Failed to delete comment")
		}
		return
	}

	// 检查是否是评论作者
	if comment.UserID != userID.(uint) {
		handlers.Forbidden(c, "You are not the author of this comment")
		return
	}

	// 删除评论（软删除）
	if err := database.DB.Delete(&comment).Error; err != nil {
		middleware.ErrorLogger(err, "Failed to delete comment")
		handlers.InternalServerError(c, "Failed to delete comment")
		return
	}

	handlers.Success(c, "Comment deleted successfully", nil)
}
