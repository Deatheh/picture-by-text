package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// SaveRefreshToken сохраняет refresh_token для пользователя
func (r *RedisClient) SaveRefreshToken(ctx context.Context, userID, token string, ttlSeconds int) error {
	return r.Client.Set(ctx, "refresh:"+userID, token, time.Duration(ttlSeconds)*time.Second).Err()
}

// GetRefreshToken возвращает сохранённый refresh_token для пользователя
func (r *RedisClient) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	return r.Client.Get(ctx, "refresh:"+userID).Result()
}

// DeleteRefreshToken удаляет refresh_token (при логауте)
func (r *RedisClient) DeleteRefreshToken(ctx context.Context, userID string) error {
	return r.Client.Del(ctx, "refresh:"+userID).Err()
}
