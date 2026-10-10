package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	Name      string
	PostCount int
	Posts     []Post
}

type Post struct {
	gorm.Model
	Title         string
	Content       string
	CommentStatus string // "有评论" / "无评论"
	UserID        uint
	Comments      []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
}

// 题目3.1：文章创建后，自动把作者的 PostCount +1
// 使用钩子传入的 tx，与创建文章处于同一事务，失败会一起回滚
func (p *Post) AfterCreate(tx *gorm.DB) error {
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error
}

// 题目3.2：评论删除后，若文章已无评论，则把文章的评论状态改为 "无评论"
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 软删除的记录会被 Count 自动排除
	var count int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			UpdateColumn("comment_status", "无评论").Error
	}
	return nil
}

func initMockData(db *gorm.DB) {
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

func printPostsWithCommentsOfUser(db *gorm.DB, name string) {
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

func printMostCommentedPost(db *gorm.DB) {
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

func publishPost(db *gorm.DB) {
	var user User
	if err := db.Where("name = ?", "Bob").First(&user).Error; err != nil {
		log.Printf("查询用户失败: %v\n", err)
		return
	}
	fmt.Printf("📮 发文前 %s 的文章数: %d\n", user.Name, user.PostCount)

	// 通过 UserID 关联创建文章，Create 成功后会触发 Post 的 AfterCreate 钩子
	post := Post{
		Title:   "钩子函数测试文章",
		Content: "验证 AfterCreate 是否自动更新 PostCount",
		UserID:  user.ID,
	}
	if err := db.Create(&post).Error; err != nil {
		log.Printf("发布文章失败: %v\n", err)
		return
	}

	// 重新查询，拿到钩子更新后的数据库值
	var updated User
	if err := db.First(&updated, user.ID).Error; err != nil {
		log.Printf("重新查询用户失败: %v\n", err)
		return
	}
	fmt.Printf("📮 发文后 %s 的文章数: %d\n", updated.Name, updated.PostCount)
}

func deleteComment(db *gorm.DB) {
	// 找到只有 1 条评论的文章（"Go 语言并发指南"）
	var post Post
	err := db.Where("title = ?", "Go 语言并发指南").First(&post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("没有找到该文章")
		return
	}
	if err != nil {
		log.Printf("查询文章失败: %v\n", err)
		return
	}

	// 必须先查出完整的 Comment 对象再删除，钩子里才能拿到 PostID
	var comment Comment
	if err := db.Where("post_id = ?", post.ID).First(&comment).Error; err != nil {
		fmt.Println("🗑  该文章已没有评论可删，跳过")
	} else if err := db.Delete(&comment).Error; err != nil {
		log.Printf("删除评论失败: %v\n", err)
		return
	}

	var updated Post
	db.First(&updated, post.ID)
	fmt.Printf("🗑  文章 [%s] 的评论状态: %q\n", updated.Title, updated.CommentStatus)
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
	initMockData(db)

	// 应题目2.1
	printPostsWithCommentsOfUser(db, "Alice")
	printPostsWithCommentsOfUser(db, "Bob")

	// 应题目2.2
	printMostCommentedPost(db)

	// 应题目3.1：创建文章触发 AfterCreate
	publishPost(db)

	// 应题目3.2：删除评论触发 AfterDelete
	deleteComment(db)
}
