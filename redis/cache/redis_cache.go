package cache

import (
	"github.com/redis/go-redis/v9"
)

type RedisCacher struct {
	client *redis.Client
}

func NewRedisCacher(opt redis.Options) *RedisCacher {
	return &RedisCacher{
		client: redis.NewClient(&redis.Options{
			Addr: opt.Addr,
		}),
	}
}
