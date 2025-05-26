package storage

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// User 数据库模型
type User struct {
	ID                int64
	UserID            string
	Username          string
	PasswordHash      string
	UserLike          string
	UserLikeEmbedding []float32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// PostgresStorage 封装数据库操作
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage 创建数据库连接
func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

// CreateUser 创建用户
func (s *PostgresStorage) CreateUser(user *User) error {
	query := `
        INSERT INTO users (user_id, username, password_hash, user_like, user_like_embedding, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := s.db.Exec(query, user.UserID, user.Username, user.PasswordHash, user.UserLike, float32SliceToPGVector(user.UserLikeEmbedding), user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByUsername 根据用户名查找用户
func (s *PostgresStorage) GetUserByUsername(username string) (*User, error) {
	query := `
        SELECT id, user_id, username, password_hash, user_like, user_like_embedding, created_at, updated_at
        FROM users WHERE username = $1
    `
	row := s.db.QueryRow(query, username)
	return scanUser(row)
}

// GetUserByUserID 根据user_id查找用户
func (s *PostgresStorage) GetUserByUserID(userID string) (*User, error) {
	query := `
        SELECT id, user_id, username, password_hash, user_like, user_like_embedding, created_at, updated_at
        FROM users WHERE user_id = $1
    `
	row := s.db.QueryRow(query, userID)
	return scanUser(row)
}

// scanUser 从sql.Row解析User
func scanUser(row *sql.Row) (*User, error) {
	var u User
	var likeEmbedding []byte // 用于接收向量类型
	err := row.Scan(&u.ID, &u.UserID, &u.Username, &u.PasswordHash, &u.UserLike, &likeEmbedding, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	u.UserLikeEmbedding, err = pgVectorToFloat32Slice(likeEmbedding)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// float32SliceToPGVector 将float32切片转为PostgreSQL向量类型字符串
func float32SliceToPGVector(floats []float32) string {
	s := "["
	for i, f := range floats {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("%f", f)
	}
	s += "]"
	return s
}

// pgVectorToFloat32Slice 将PostgreSQL向量类型转为float32切片
func pgVectorToFloat32Slice(b []byte) ([]float32, error) {
	// 假设数据库返回格式为 {0.1,0.2,0.3}
	str := string(b)
	str = str[1 : len(str)-1] // 去掉{}
	if str == "" {
		return []float32{}, nil
	}
	parts := strings.Split(str, ",")
	res := make([]float32, len(parts))
	for i, p := range parts {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 32)
		if err != nil {
			return nil, err
		}
		res[i] = float32(v)
	}
	return res, nil
}
