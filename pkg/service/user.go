package service

import (
	"context"
	"errors"
	"math/rand/v2"
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
	span.SetTag("user_id", req.RequestId)

	if req.RequestId == "" {
		return nil, errors.New("request_id is required")
	}
	// 幂等：先查用户名
	exist, _ := s.Store.GetUserByUserID(req.RequestId)
	if exist != nil {
		return nil, errors.New("username already exists")
	}
	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// 获取 user_like_embedding（此处用假数据，实际应调用外部服务）
	userLikeEmbedding := make([]float32, 768)
	for i := range userLikeEmbedding {
		userLikeEmbedding[i] = rand.Float32()
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
	return &pb.RegisterResponse{UserId: userID}, nil
}

// Login 用户登录，返回JWT
func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	span := opentracing.StartSpan("UserService.Login")
	defer span.Finish()
	span.SetTag("username", req.Username)

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
