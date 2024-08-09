package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ZenSam7/Education/tools"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"time"
)

type RedisCacher struct {
	client        *redis.Client
	CacheDuration time.Duration
}

func NewRedisCacher(opt redis.Options, config tools.Config) *RedisCacher {
	rC := &RedisCacher{
		CacheDuration: config.CacheDuration,
		client: redis.NewClient(&redis.Options{
			Addr: opt.Addr,
		}),
	}

	// Проверка соединения с Redis
	_, err := rC.client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("не получилось создать redis cacher")
	}

	return rC
}

func (r *RedisCacher) GetCache(ctx context.Context, key string, dest any) error {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return errors.New("нету кеша")
	} else if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCacher) SetCache(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.CacheDuration).Err()
}
