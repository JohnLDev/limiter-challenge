package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/johnldev/rate-limiter/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	conn *redis.Client
	ctx  context.Context
	conf config.Config
}

func (r *RedisRepository) Count(key string) (int, error) {
	iter := r.conn.Scan(r.ctx, 0, fmt.Sprintf("request:%s:*", key), 0).Iterator()
	count := 0
	for iter.Next(r.ctx) {
		count++
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	fmt.Println(count, "count")
	return count, nil
}

func (r *RedisRepository) Save(key string, id string) error {
	redisKey := fmt.Sprintf("request:%s:%s", key, id)
	fmt.Println(redisKey)
	err := r.conn.Set(r.ctx, redisKey, id, time.Second).Err()
	return err
}

func (r *RedisRepository) CheckLock(key string) (bool, error) {
	redisKey := fmt.Sprintf("lock:%s", key)
	_, err := r.conn.Get(r.ctx, redisKey).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisRepository) LockKey(key string) error {
	redisKey := fmt.Sprintf("lock:%s", key)

	err := r.conn.Set(r.ctx, redisKey, redisKey, time.Millisecond*time.Duration(r.conf.BlockTime)).Err()
	return err

}

func NewRedisRepository(ctx context.Context, config config.Config) *RedisRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.DbHost, config.DbPort),
		Password: config.DbPass,
		DB:       0,
	})

	return &RedisRepository{conn: rdb, ctx: ctx, conf: config}
}
