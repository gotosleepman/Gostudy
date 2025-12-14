package main

import (
	"fmt"
	"gorm.io/gorm"
)

// ============================================
// 钩子函数实现示例
// ============================================

// AfterCreate Post模型的钩子函数：在文章创建后自动更新用户的文章数量
// 注意：需要在Post模型中添加此方法
func (p *Post) AfterCreate(tx *gorm.DB) error {
	// 更新用户的文章数量统计
	result := tx.Model(&User{}).
		Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + ?", 1))
	
	if result.Error != nil {
		return fmt.Errorf("更新用户文章数量失败: %v", result.Error)
	}
	
	fmt.Printf("✓ 文章创建成功，已更新用户ID=%d的文章数量\n", p.UserID)
	return nil
}

// AfterDelete Comment模型的钩子函数：在评论删除后检查文章的评论数量
// 注意：需要在Comment模型中添加此方法
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 统计该文章的剩余评论数量
	var commentCount int64
	result := tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&commentCount)
	
	if result.Error != nil {
		return fmt.Errorf("统计评论数量失败: %v", result.Error)
	}
	
	// 如果评论数量为0，更新文章的评论状态为"无评论"
	if commentCount == 0 {
		updateResult := tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			Update("comment_status", "无评论")
		
		if updateResult.Error != nil {
			return fmt.Errorf("更新文章评论状态失败: %v", updateResult.Error)
		}
		
		fmt.Printf("✓ 评论删除成功，文章ID=%d的评论数量为0，已更新状态为'无评论'\n", c.PostID)
	} else {
		fmt.Printf("✓ 评论删除成功，文章ID=%d还有%d条评论\n", c.PostID, commentCount)
	}
	
	return nil
}

// ============================================
// 模型定义（需要添加的字段）
// ============================================

/*
// User 用户模型 - 需要添加 PostCount 字段
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string    `gorm:"not null;size:255" json:"-"`
	PostCount int       `gorm:"default:0" json:"post_count"` // 新增：文章数量统计字段
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Posts     []Post    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"posts,omitempty"`
}

// Post 文章模型 - 需要添加 CommentStatus 字段和 AfterCreate 钩子函数
type Post struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"not null;size:255" json:"title"`
	Content       string    `gorm:"type:text" json:"content"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	CommentStatus string    `gorm:"default:'有评论';size:50" json:"comment_status"` // 新增：评论状态字段
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Comments      []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}

// Comment 评论模型 - 需要添加 AfterDelete 钩子函数
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
*/

