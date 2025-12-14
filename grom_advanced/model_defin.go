package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string    `gorm:"not null;size:255" json:"-"` // json:"-" 表示序列化时隐藏密码
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 一对多关系：一个用户可以有多篇文章
	Posts []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"posts,omitempty"`
}

// Post 文章模型
type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null;size:255" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	UserID    uint      `gorm:"not null;index" json:"user_id"` // 外键
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

func main() {
	// 数据库连接配置
	// 请根据实际情况修改数据库连接信息
	// 格式: username:password@tcp(host:port)/database_name?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "root:psaaword@tcp(localhost:3306)/gostudy?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	fmt.Println("数据库连接成功！")

	// 自动迁移：创建或更新数据库表结构
	// 这会根据模型定义自动创建表，如果表已存在则更新结构
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	fmt.Println("数据库表创建成功！")
	fmt.Println("- users 表已创建")
	fmt.Println("- posts 表已创建")
	fmt.Println("- comments 表已创建")

	// 可选：创建一些示例数据
	createSampleData(db)
}

// createSampleData 创建示例数据（可选）
func createSampleData(db *gorm.DB) {
	// 检查是否已有数据
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount > 0 {
		fmt.Println("数据库中已有数据，跳过示例数据创建")
		return
	}

	// 创建用户
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password_here",
	}
	db.Create(&user)
	fmt.Printf("创建用户: %s (ID: %d)\n", user.Username, user.ID)

	// 创建文章
	post := Post{
		Title:   "我的第一篇文章",
		Content: "这是文章的内容...",
		UserID:  user.ID,
	}
	db.Create(&post)
	fmt.Printf("创建文章: %s (ID: %d)\n", post.Title, post.ID)

	// 创建评论
	comment := Comment{
		Content: "这是一条评论",
		PostID:  post.ID,
		UserID:  user.ID,
	}
	db.Create(&comment)
	fmt.Printf("创建评论 (ID: %d)\n", comment.ID)

	fmt.Println("示例数据创建完成！")
}
