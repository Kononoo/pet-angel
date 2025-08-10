package main

import (
	"context"
	"fmt"
	"log"

	"pet-angel/internal/conf"
	"pet-angel/internal/data"
)

func testDataQueries() {
	// 加载配置
	c := &conf.Bootstrap{
		Data: &conf.Data{
			Database: &conf.Data_Database{
				Source: "root:140001@tcp(47.121.139.174:3306)/pet_angel?parseTime=True&loc=Local",
			},
		},
	}

	// 初始化数据层
	d, cleanup, err := data.NewData(c.Data, nil)
	if err != nil {
		log.Fatalf("failed to init data: %v", err)
	}
	defer cleanup()

	ctx := context.Background()

	// 测试社区数据查询
	fmt.Println("=== 测试社区数据查询 ===")
	communityRepo := data.NewCommunityRepo(d)

	// 测试分类列表
	categories, err := communityRepo.ListCategories(ctx)
	if err != nil {
		fmt.Printf("ListCategories error: %v\n", err)
	} else {
		fmt.Printf("Categories count: %d\n", len(categories))
		for _, cat := range categories {
			fmt.Printf("  - %s (ID: %d)\n", cat.Name, cat.ID)
		}
	}

	// 测试帖子列表
	total, posts, err := communityRepo.ListPosts(ctx, 1, 0, 0, "", 1, 10)
	if err != nil {
		fmt.Printf("ListPosts error: %v\n", err)
	} else {
		fmt.Printf("Posts total: %d, current page: %d\n", total, len(posts))
		for _, post := range posts {
			fmt.Printf("  - %s (ID: %d, User: %s)\n", post.Title, post.ID, post.User.Nickname)
		}
	}

	// 测试消息数据查询
	fmt.Println("\n=== 测试消息数据查询 ===")
	messageRepo := data.NewMessageRepo(d)

	// 测试消息列表
	msgTotal, messages, err := messageRepo.ListMessages(ctx, 1, false, 1, 10)
	if err != nil {
		fmt.Printf("ListMessages error: %v\n", err)
	} else {
		fmt.Printf("Messages total: %d, current page: %d\n", msgTotal, len(messages))
		for _, msg := range messages {
			sender := "用户"
			if msg.Sender == 1 {
				sender = "AI"
			}
			fmt.Printf("  - %s: %s (ID: %d)\n", sender, msg.Content, msg.ID)
		}
	}

	// 测试小纸条列表
	noteTotal, notes, err := messageRepo.ListMessages(ctx, 1, true, 1, 10)
	if err != nil {
		fmt.Printf("ListNotes error: %v\n", err)
	} else {
		fmt.Printf("Notes total: %d, current page: %d\n", noteTotal, len(notes))
		for _, note := range notes {
			status := "已解锁"
			if note.IsLocked {
				status = fmt.Sprintf("锁定(%d金币)", note.UnlockCoins)
			}
			fmt.Printf("  - %s [%s] (ID: %d)\n", note.Content, status, note.ID)
		}
	}

	// 测试用户数据查询
	fmt.Println("\n=== 测试用户数据查询 ===")
	userRepo := data.NewUserRepo(d)

	// 测试用户信息
	user, err := userRepo.GetUserByID(ctx, 1)
	if err != nil {
		fmt.Printf("GetUserByID error: %v\n", err)
	} else {
		fmt.Printf("User: %s (ID: %d, Coins: %d)\n", user.Nickname, user.Id, user.Coins)
	}

	// 测试用户帖子
	userPosts, err := userRepo.GetUserPostsBrief(ctx, 1, 5)
	if err != nil {
		fmt.Printf("GetUserPostsBrief error: %v\n", err)
	} else {
		fmt.Printf("User posts count: %d\n", len(userPosts))
		for _, post := range userPosts {
			fmt.Printf("  - %s (ID: %d, Likes: %d)\n", post.Title, post.Id, post.LikedCount)
		}
	}

	// 测试关注列表
	follows, followTotal, err := userRepo.GetFollowList(ctx, 1, 1, 10)
	if err != nil {
		fmt.Printf("GetFollowList error: %v\n", err)
	} else {
		fmt.Printf("Follows total: %d, current page: %d\n", followTotal, len(follows))
		for _, follow := range follows {
			fmt.Printf("  - %s (ID: %d)\n", follow.Nickname, follow.Id)
		}
	}

	// 测试认证数据查询
	fmt.Println("\n=== 测试认证数据查询 ===")
	authRepo := data.NewAuthRepo(d)

	// 测试用户登录
	authUser, err := authRepo.GetByUsername(ctx, "lemon")
	if err != nil {
		fmt.Printf("GetByUsername error: %v\n", err)
	} else {
		fmt.Printf("Auth user: %s (ID: %d, Coins: %d)\n", authUser.Nickname, authUser.Id, authUser.Coins)
	}

	fmt.Println("\n=== 数据查询测试完成 ===")
}
