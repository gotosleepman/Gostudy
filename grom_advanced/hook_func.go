package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户模型（添加文章数量统计字段）
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string    `gorm:"not null;size:255" json:"-"`  // json:"-" 表示序列化时隐藏密码
	PostCount int       `gorm:"default:0" json:"post_count"` // 文章数量统计字段
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 一对多关系：一个用户可以有多篇文章
	Posts []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"posts,omitempty"`
}

// Post 文章模型（添加评论状态字段）
type Post struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"not null;size:255" json:"title"`
	Content       string    `gorm:"type:text" json:"content"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`               // 外键
	CommentStatus string    `gorm:"default:'有评论';size:50" json:"comment_status"` // 评论状态：有评论/无评论
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	// 多对一关系：一篇文章属于一个用户
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// 一对多关系：一篇文章可以有多个评论
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}

// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	PostID    uint      `gorm:"not null;index" json:"post_id"` // 外键
	UserID    uint      `gorm:"not null;index" json:"user_id"` // 可选：评论也可以关联到用户
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 多对一关系：一个评论属于一篇文章
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
	// 多对一关系：一个评论属于一个用户（可选）
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// AfterCreate Post模型的钩子函数：在文章创建后自动更新用户的文章数量
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

func main() {
	// 数据库连接配置
	dsn := "root:123456@tcp(localhost:3306)/gostudy?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	fmt.Println("数据库连接成功！")

	// 自动迁移：创建或更新数据库表结构
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	fmt.Println("数据库表创建/更新成功！")
	fmt.Println("- users 表已更新（添加post_count字段）")
	fmt.Println("- posts 表已更新（添加comment_status字段）")
	fmt.Println("- comments 表已创建")

	// 测试钩子函数
	testHooks(db)
}

// testHooks 测试钩子函数
func testHooks(db *gorm.DB) {
	fmt.Println("\n=== 测试钩子函数 ===")

	// 1. 检查是否有用户，如果没有则创建一个
	var user User
	result := db.First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		user = User{
			Username:  "hooktest",
			Email:     "hooktest@example.com",
			Password:  "password123",
			PostCount: 0,
		}
		db.Create(&user)
		fmt.Printf("创建测试用户: %s (ID: %d, 文章数量: %d)\n", user.Username, user.ID, user.PostCount)
	} else {
		fmt.Printf("使用现有用户: %s (ID: %d, 文章数量: %d)\n", user.Username, user.ID, user.PostCount)
	}

	// 2. 测试Post的AfterCreate钩子：创建文章
	fmt.Println("\n--- 测试1: 创建文章，触发AfterCreate钩子 ---")
	post := Post{
		Title:         "测试文章 - 钩子函数",
		Content:       "这是一篇用于测试钩子函数的文章",
		UserID:        user.ID,
		CommentStatus: "有评论",
	}
	db.Create(&post)
	fmt.Printf("文章创建成功: %s (ID: %d)\n", post.Title, post.ID)

	// 查询用户，验证文章数量是否更新
	db.First(&user, user.ID)
	fmt.Printf("用户文章数量已更新为: %d\n", user.PostCount)

	// 3. 测试Comment的AfterDelete钩子：创建评论然后删除
	fmt.Println("\n--- 测试2: 创建评论 ---")
	comment := Comment{
		Content: "这是一条测试评论",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	db.Create(&comment)
	fmt.Printf("评论创建成功 (ID: %d)\n", comment.ID)

	// 查询文章，查看评论状态
	var postCheck Post
	db.First(&postCheck, post.ID)
	fmt.Printf("文章评论状态: %s\n", postCheck.CommentStatus)

	// 删除评论，触发AfterDelete钩子
	fmt.Println("\n--- 测试3: 删除评论，触发AfterDelete钩子 ---")
	db.Delete(&comment)
	fmt.Printf("评论已删除 (ID: %d)\n", comment.ID)

	// 再次查询文章，验证评论状态是否更新
	db.First(&postCheck, post.ID)
	fmt.Printf("文章评论状态已更新为: %s\n", postCheck.CommentStatus)

	fmt.Println("\n=== 钩子函数测试完成 ===")
}
