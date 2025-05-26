package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/golang-jwt/jwt/v4"

	pb "github.com/LtePrince/GO-grpc-server/pkg/api"
	"github.com/LtePrince/GO-grpc-server/pkg/service"
	"github.com/LtePrince/GO-grpc-server/pkg/storage"
)

// JWT拦截器，校验token并将user_id注入metadata
func AuthInterceptor(secret []byte) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// 跳过无需鉴权的方法
		if info.FullMethod == "/user.UserService/Register" || info.FullMethod == "/user.UserService/Login" {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata")
		}
		auths := md.Get("authorization")
		if len(auths) == 0 {
			return nil, fmt.Errorf("missing authorization token")
		}
		tokenStr := auths[0]
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			return nil, fmt.Errorf("invalid token: %v", err)
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}
		// 将user_id注入metadata
		newMD := metadata.Pairs("user_id", userID)
		ctx = metadata.NewIncomingContext(ctx, metadata.Join(md, newMD))
		return handler(ctx, req)
	}
}

// 提供PostgresStorage
func NewPostgresStorage() (*storage.PostgresStorage, error) {
	dsn := "host=localhost port=5432 user=adolph password=Wzy20031003 dbname=grpctest sslmode=disable"
	return storage.NewPostgresStorage(dsn)
}

// 提供UserServiceServer
func NewUserServiceServer(store *storage.PostgresStorage) *service.UserServiceServer {
	jwtSecret := "your_jwt_secret"
	return service.NewUserServiceServer(store, jwtSecret)
}

// 提供SystemServiceServer
func NewSystemServiceServer() *service.SystemServiceServer {
	return service.NewSystemServiceServer("./data")
}

// 提供gRPC Server
func NewGRPCServer(userSrv *service.UserServiceServer) *grpc.Server {
	jwtSecret := userSrv.JWTSecret
	return grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor(jwtSecret)),
	)
}

// 注册服务并启动
func RegisterServer(lc fx.Lifecycle, srv *grpc.Server, userSrv *service.UserServiceServer, sysSrv *service.SystemServiceServer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":50051")
			if err != nil {
				return err
			}
			pb.RegisterUserServiceServer(srv, userSrv)
			pb.RegisterSystemServiceServer(srv, sysSrv)
			go func() {
				log.Println("gRPC server listening on :50051")
				if err := srv.Serve(lis); err != nil {
					log.Fatalf("failed to serve: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.GracefulStop()
			return nil
		},
	})
}

func main() {
	fx.New(
		fx.Provide(
			NewPostgresStorage,
			NewUserServiceServer,
			NewSystemServiceServer,
			NewGRPCServer,
		),
		fx.Invoke(RegisterServer),
	).Run()
}
