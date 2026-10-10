package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	Name  string
	Posts []Post
}

type Post struct {
	gorm.Model
	Title    string
	Content  string
	UserID   uint
	Comments []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
}

func init_mock_data(db *gorm.DB) {
	// 每次启动先检查是否已有数据，防止重复插入
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		// 构造嵌套测试数据 (Mock Data)
		users := []User{
			{
				Name: "Alice",
				Posts: []Post{
					{
						Title:   "Go 语言并发指南",
						Content: "Goroutine 和 Channel 太好用了！",
						Comments: []Comment{
							{Content: "写得好，沙发！"},
						},
					},
					{
						Title:   "GORM 踩坑实录",
						Content: "今天学习了自动迁移和级联创建...",
						Comments: []Comment{ // 这篇刻意造 3 条评论，用于测试“最多评论”
							{Content: "催更！"},
							{Content: "学到了学到了。"},
							{Content: "楼主牛逼！"},
						},
					},
				},
			},
			{
				Name: "Bob",
				Posts: []Post{
					{
						Title:   "深夜碎碎念",
						Content: "今天真是充实的一天",
					},
				},
			},
		}

		// 魔法发生在这里：一次性插入所有用户、文章和评论！
		if err := db.Create(&users).Error; err != nil {
			log.Fatal("插入测试数据失败:", err)
		}
		fmt.Println("🌱 测试数据 (Mock Data) 插入成功！")
	} else {
		fmt.Println("✅ 数据库中已有样本数据，跳过插入。")
	}
}

func print_post_with_comment_of_user(db *gorm.DB, name string) {
	var user User

	// Preload("Posts") 查出文章，Preload("Posts.Comments") 顺藤摸瓜查出文章下的评论
	// 这里会做基于依赖路径的递归推导
	err := db.Preload("Posts.Comments").Where("name = ?", name).First(&user).Error
	if err != nil {
		log.Printf("查询用户 %s 失败: %v\n", name, err)
		return
	}

	fmt.Println("=====================================")
	fmt.Printf("👤 用户: %s\n", user.Name)
	for _, post := range user.Posts {
		fmt.Printf("  📝 文章: [%s] %s\n", post.Title, post.Content)
		for _, comment := range post.Comments {
			fmt.Printf("      💬 评论: %s\n", comment.Content)
		}
		if len(post.Comments) == 0 {
			fmt.Println("      (暂无评论)")
		}
	}
	fmt.Println("=====================================")
}

// 引入 DTO (Data Transfer Object) 隔离视图数据与底层模型
type PostWithCount struct {
	Post             // 匿名嵌套：继承 Post 原有的所有字段
	CommentCount int // 专门接收 COUNT 出来的别名数据
}

func print_most_commented_post(db *gorm.DB) {
	var result PostWithCount // 使用 DTO 作为结果接收器

	err := db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count DESC").
		First(&result).Error

	if err != nil {
		log.Println("查询失败:", err)
		return
	}

	fmt.Println("=====================================")
	fmt.Printf("🔥 评论最多的文章揭晓！\n")
	fmt.Printf("📝 标题: [%s]\n", result.Title)
	fmt.Printf("💬 内容: %s\n", result.Content)
	fmt.Printf("📊 统计: 共计 %d 条评论\n", result.CommentCount)
	fmt.Println("=====================================")
}

func main() {
	// 应题目要求1
	// 连接本地 SQLite 数据库（如果文件不存在会自动创建）
	db, err := gorm.Open(sqlite.Open("data/blog.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	if err := db.AutoMigrate(&User{}, &Post{}, &Comment{}); err != nil {
		log.Fatal("failed to AutoMigrate database")
	}

	fmt.Println("hello GORM")

	//初始化填充测数据
	init_mock_data(db)

	// 应题目2.1
	print_post_with_comment_of_user(db, "Alice")
	print_post_with_comment_of_user(db, "Bob")

	// 应题目2.2
	print_most_commented_post(db)

}
