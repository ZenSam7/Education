package cache

import (
	"context"
)

type Cacher interface {
	GetCache(ctx context.Context, key string, dest any) error
	SetCache(ctx context.Context, key string, value any) error
}
