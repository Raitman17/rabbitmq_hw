package storage

import (
	"context"
	"os"
	"time"
	"github.com/redis/go-redis/v9"
)

type AbsRedis interface {
	ConnectToRedis()
	Set(key string, value string)
	Get(key string) *string
	Close()
}

type Redis struct {
	rdb *redis.Client
}

func NewRedis() AbsRedis {
	return &Redis{}
}

func (r *Redis) ConnectToRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		failOnError(err, "Failed to connect to Redis")
	}

	r.rdb = rdb
}

func (r *Redis) Set(key string, value string) {
	err := r.rdb.Set(context.Background(), key, value, 10*time.Second).Err()
	if err != nil {
        failOnError(err, "Failed to add record to redis")
    }
}

func (r *Redis) Get(key string) *string {
	val, err := r.rdb.Get(context.Background(), key).Result()
	switch {
	case err == redis.Nil:
		return nil
	case err != nil:
		failOnError(err, "Failed to get message from cache")
	}
	return &val
}

func (r *Redis) Close() {
	r.rdb.Close()
}
