package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/LtePrince/GO-grpc-server/pkg/api"
)

func main() {
	// 推荐新写法
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	// 1. 注册
	registerResp, err := client.Register(context.Background(), &pb.RegisterRequest{
		Username:  "testuser",
		Password:  "testpass",
		Like:      "golang",
		RequestId: fmt.Sprintf("req-%d", time.Now().UnixNano()),
	})
	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}
	fmt.Println("Register user_id:", registerResp.UserId)

	// 2. 登录
	loginResp, err := client.Login(context.Background(), &pb.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Println("Login access_token:", loginResp.AccessToken)

	// 3. 获取用户信息（带token）
	md := metadata.New(map[string]string{"authorization": loginResp.AccessToken})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	userInfo, err := client.GetUserInfo(ctx, &pb.GetUserInfoRequest{})
	if err != nil {
		log.Fatalf("GetUserInfo failed: %v", err)
	}
	fmt.Printf("UserInfo: user_id=%s, username=%s, like=%s, embedding=%v, create_at=%s\n",
		userInfo.UserId, userInfo.Username, userInfo.Like,
		func() []float32 {
			if len(userInfo.LikeEmbedding) > 5 {
				return userInfo.LikeEmbedding[:5]
			}
			return userInfo.LikeEmbedding
		}(),
		userInfo.CreateAt)
}
