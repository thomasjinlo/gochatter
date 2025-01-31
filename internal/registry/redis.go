package registry

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisRegistry struct {
	client *redis.Client
}

func New(addr, password string, db int) *RedisRegistry {
	return &RedisRegistry{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *RedisRegistry) Get(ctx context.Context, clientId, endpoint string) ([]string, error) {
	return r.client.SMembers(ctx, clientId).Result()
}

func (r *RedisRegistry) Add(ctx context.Context, clientId, endpoint string) error {
	return r.client.SAdd(ctx, clientId, endpoint).Err()
}

func (r *RedisRegistry) Remove(ctx context.Context, clientId, endpoint string) error {
	return r.client.SRem(ctx, clientId, endpoint).Err()
}
