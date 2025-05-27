package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/metadata"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	pb "github.com/LtePrince/GO-grpc-server/pkg/api"
	"github.com/LtePrince/GO-grpc-server/pkg/storage"
)

// UserServiceServer 实现 pb.UserServiceServer
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	Store     *storage.PostgresStorage
	JWTSecret []byte
}

// NewUserServiceServer 构造函数
func NewUserServiceServer(store *storage.PostgresStorage, secret string) *UserServiceServer {
	return &UserServiceServer{
		Store:     store,
		JWTSecret: []byte(secret),
	}
}

// Register 用户注册（幂等）
func (s *UserServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	span := opentracing.StartSpan("UserService.Register")
	defer span.Finish()

	if req.RequestId == "" {
		return nil, errors.New("request_id is required")
	}
	// 幂等：先查用户名
	exist, _ := s.Store.GetUserByUsername(req.Username)
	if exist != nil {
		return nil, errors.New("username already exists")
	}
	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// 获取 user_like_embedding（此处用假数据，实际应调用外部服务）
	// 调用 HuggingFace API 获取 embedding
	userLikeEmbedding, err := fetchEmbeddingFromHuggingFace(req.Like)
	if err != nil {
		return nil, fmt.Errorf("embedding service error: %v", err)
	}
	now := time.Now()
	userID := uuid.NewString()
	user := &storage.User{
		UserID:            userID,
		Username:          req.Username,
		PasswordHash:      string(hash),
		UserLike:          req.Like,
		UserLikeEmbedding: userLikeEmbedding,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	err = s.Store.CreateUser(user)
	if err != nil {
		return nil, err
	}
	span.SetTag("user_id", userID)
	return &pb.RegisterResponse{UserId: userID}, nil
}

// Login 用户登录，返回JWT
func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	span := opentracing.StartSpan("UserService.Login")
	defer span.Finish()

	user, err := s.Store.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}
	// 生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(s.JWTSecret)
	if err != nil {
		return nil, err
	}
	span.SetTag("username", user.UserID)
	return &pb.LoginResponse{AccessToken: tokenString}, nil
}

// GetUserInfo 获取用户信息（需鉴权）
func (s *UserServiceServer) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	userID, err := getUserIDFromContext(ctx, s.JWTSecret)
	span := opentracing.StartSpan("UserService.GetUserInfo")
	defer span.Finish()
	span.SetTag("user_id", userID)

	if err != nil {
		return nil, err
	}

	user, err := s.Store.GetUserByUserID(userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}
	return &pb.GetUserInfoResponse{
		UserId:        user.UserID,
		Username:      user.Username,
		Like:          user.UserLike,
		LikeEmbedding: user.UserLikeEmbedding,
		CreateAt:      user.CreatedAt.Format(time.RFC3339),
		UpdateAt:      user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// getUserIDFromContext 从 context 的 metadata 里解析 JWT 并获取 user_id
func getUserIDFromContext(ctx context.Context, secret []byte) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}
	auths := md.Get("authorization")
	if len(auths) == 0 {
		return "", errors.New("missing authorization token")
	}
	tokenStr := auths[0]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return "", errors.New("invalid token")
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	return userID, nil
}

// 调用 HuggingFace API 获取 embedding
func fetchEmbeddingFromHuggingFace(text string) ([]float32, error) {
	keyData, err := os.ReadFile("./key/hf_key.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read huggingface api key: %v", err)
	}
	apiToken := string(keyData)
	url := "https://router.huggingface.co/hf-inference/models/intfloat/multilingual-e5-large-instruct/pipeline/feature-extraction"

	requestData := map[string]string{"inputs": text}
	jsonData, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("huggingface api error: %s, body: %s", resp.Status, string(body))
	}

	// 先尝试二维数组
	var response2 [][]float64
	if err := json.Unmarshal(body, &response2); err == nil && len(response2) > 0 {
		embedding := make([]float32, len(response2[0]))
		for i, v := range response2[0] {
			embedding[i] = float32(v)
		}
		return embedding, nil
	}

	// 再尝试一维数组
	var response1 []float64
	if err := json.Unmarshal(body, &response1); err == nil && len(response1) > 0 {
		embedding := make([]float32, len(response1))
		for i, v := range response1 {
			embedding[i] = float32(v)
		}
		return embedding, nil
	}

	return nil, fmt.Errorf("embedding response format error, body: %s", string(body))
}
