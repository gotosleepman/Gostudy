package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	PostID    uint           `gorm:"not null" json:"post_id"`
	Post      Post           `gorm:"foreignKey:PostID" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 创建评论请求
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
	PostID  uint   `json:"post_id" binding:"required"`
}
