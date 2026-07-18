package cache

import (
	"context"
	"time"
)

type CacheGetRequest interface {
	Do(ctx context.Context) error
}

type CacheSetRequest interface {
	WithTTL(ttl time.Duration) CacheSetRequest
	WithExpiry(expiry time.Time) CacheSetRequest
	Do(ctx context.Context) error
}

type Cache interface {
	Get(namespace, key string, response interface{}) CacheGetRequest
	Set(namespace, key string, value interface{}) CacheSetRequest
}
