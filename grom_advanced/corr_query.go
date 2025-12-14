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

// QueryUserPostsWithComments 查询某个用户发布的所有文章及其对应的评论信息
func QueryUserPostsWithComments(db *gorm.DB, userID uint) {
	fmt.Printf("\n=== 查询用户ID=%d发布的所有文章及其评论 ===\n", userID)

	var user User
	// 使用Preload预加载文章和评论数据
	err := db.Preload("Posts.Comments").First(&user, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Printf("用户ID=%d不存在\n", userID)
		} else {
			log.Printf("查询失败: %v", err)
		}
		return
	}

	fmt.Printf("用户: %s (ID: %d)\n", user.Username, user.ID)
	fmt.Printf("文章数量: %d\n", len(user.Posts))

	for i, post := range user.Posts {
		fmt.Printf("\n文章 %d:\n", i+1)
		fmt.Printf("  标题: %s\n", post.Title)
		fmt.Printf("  内容: %s\n", post.Content)
		fmt.Printf("  发布时间: %s\n", post.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("  评论数量: %d\n", len(post.Comments))

		if len(post.Comments) > 0 {
			fmt.Println("  评论列表:")
			for j, comment := range post.Comments {
				fmt.Printf("    %d. %s (用户ID: %d, 时间: %s)\n",
					j+1,
					comment.Content,
					comment.UserID,
					comment.CreatedAt.Format("2006-01-02 15:04:05"))
			}
		}
	}
}

// PostWithCommentCount 用于存储文章和评论数量的结构体
type PostWithCommentCount struct {
	ID           uint
	Title        string
	Content      string
	UserID       uint
	CreatedAt    time.Time
	CommentCount int64  // 评论数量
	Username     string // 用户名
}

// QueryMostCommentedPost 查询评论数量最多的文章信息
func QueryMostCommentedPost(db *gorm.DB) {
	fmt.Println("\n=== 查询评论数量最多的文章 ===")

	var result PostWithCommentCount

	// 使用原生SQL查询，联表查询评论最多的文章
	query := `
		SELECT 
			p.id, p.title, p.content, p.user_id, p.created_at,
			u.username,
			COUNT(c.id) as comment_count
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		GROUP BY p.id, p.title, p.content, p.user_id, p.created_at, u.username
		ORDER BY comment_count DESC
		LIMIT 1
	`

	err := db.Raw(query).Scan(&result).Error
	if err != nil {
		log.Printf("查询失败: %v", err)
		return
	}

	if result.ID == 0 {
		fmt.Println("没有找到任何文章")
		return
	}

	fmt.Printf("评论最多的文章信息:\n")
	fmt.Printf("文章ID: %d\n", result.ID)
	fmt.Printf("标题: %s\n", result.Title)
	fmt.Printf("内容: %s\n", result.Content)
	fmt.Printf("作者: %s (用户ID: %d)\n", result.Username, result.UserID)
	fmt.Printf("发布时间: %s\n", result.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("评论数量: %d\n", result.CommentCount)
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

	// 执行查询
	QueryUserPostsWithComments(db, 1) // 查询用户ID为1的文章
	QueryMostCommentedPost(db)        // 查询评论最多的文章
}
