package cache

import (
	"context"
	"financial/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	redisDriver = "redis"
	file        = "file"
)

type NotFoundError struct {
}

func (n NotFoundError) Error() string {
	return "key not exists"
}

type ICache interface {
	Add(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

func NewCache(cfg config.Cache) ICache {
	var cache ICache
	switch cfg.Driver {
	case redisDriver:
		cache = NewRedis(cfg)
		break
	case file:
		break
	default:
		cache = NewRedis(cfg)
	}

	return cache
}

type Redis struct {
	client *redis.Client
	prefix string
}

func NewRedis(cfg config.Cache) ICache {
	cl := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := cl.Ping(context.TODO()); err.Err() != nil {
		log.Fatalf("error on redis : %s", err)
	}
	return &Redis{client: cl, prefix: cfg.Prefix}
}

func (r Redis) Add(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := r.client.Set(ctx, r.createKey(key), value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, r.createKey(key)).Result()
	if err == redis.Nil {
		return "", NotFoundError{}
	} else if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

func (r Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.createKey(key)).Err()
}

func (r Redis) createKey(key string) string {
	return r.prefix + ":" + key
}
