package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStorage 封装Redis操作
type RedisStorage struct {
	client *redis.Client
}

// NewRedisStorage 创建Redis连接
func NewRedisStorage(addr, password string, db int) *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStorage{client: rdb}
}

// 缓存用户信息
func (r *RedisStorage) SetUser(ctx context.Context, user *User, ttl time.Duration) error {
	key := fmt.Sprintf("user:%s", user.UserID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

// 获取用户信息（优先查缓存）
func (r *RedisStorage) GetUser(ctx context.Context, userID string) (*User, error) {
	key := fmt.Sprintf("user:%s", userID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // 缓存未命中
	} else if err != nil {
		return nil, err
	}
	var user User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// 删除用户缓存
func (r *RedisStorage) DelUser(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:%s", userID)
	return r.client.Del(ctx, key).Err()
}
