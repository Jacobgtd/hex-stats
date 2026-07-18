package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Jacobgtd/hex-stats/backend/internal/configpack"
)

type LocalCacheConfig struct {
	Directory string
}

func LoadLocalCacheConfig() (LocalCacheConfig, error) {

	err := configpack.Load("localcache.config")
	if err != nil {
		return LocalCacheConfig{}, err
	}

	dir, err := configpack.String("CACHE_DIR")

	return LocalCacheConfig{
		Directory: dir,
	}, nil
}

type LocalCache struct {
	config LocalCacheConfig
	mutex  sync.Mutex
}

func NewLocalCache(config LocalCacheConfig) *LocalCache {
	return &LocalCache{
		config: config,
	}
}

type cacheEntry struct {
	Expiry    time.Time   `json:"expiry"`
	CreatedAt time.Time   `json:"created_at"`
	Payload   interface{} `json:"payload"`
}

type LocalCacheGetRequest struct {
	c         *LocalCache
	namespace string
	key       string
	response  interface{}
}

type LocalCacheSetRequest struct {
	c         *LocalCache
	namespace string
	key       string
	value     interface{}
	expiry    *time.Time
}

func (c *LocalCache) Get(namespace, key string, response interface{}) CacheGetRequest {
	return &LocalCacheGetRequest{
		c:         c,
		namespace: namespace,
		key:       key,
		response:  response,
	}

}

func (r *LocalCacheGetRequest) Do(ctx context.Context) error {
	r.c.mutex.Lock()
	defer r.c.mutex.Unlock()

	filename := filepath.Join(r.c.config.Directory, fmt.Sprintf("%s-%s.json", r.namespace, r.key))

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return err
	}

	// Remove expired entries
	if time.Now().After(entry.Expiry) {
		_ = os.Remove(filename)
		return errors.New("cache key expired")
	}

	bytes, err := json.Marshal(entry.Payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, r.response)
}

func (c *LocalCache) Set(namespace, key string, value interface{}) CacheSetRequest {
	return &LocalCacheSetRequest{
		c:         c,
		namespace: namespace,
		key:       key,
		value:     value,
		expiry:    nil,
	}
}

func (r *LocalCacheSetRequest) WithExpiry(expiry time.Time) CacheSetRequest {
	r.expiry = &expiry
	return r
}

func (r *LocalCacheSetRequest) WithTTL(ttl time.Duration) CacheSetRequest {
	expiry := time.Now().Add(ttl)
	return r.WithExpiry(expiry)
}

func (r *LocalCacheSetRequest) Do(ctx context.Context) error {

	if r.expiry == nil {
		return errors.New("expiry cannot be nil")
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(r.c.config.Directory, 0755); err != nil {
		return err
	}

	filename := filepath.Join(r.c.config.Directory, fmt.Sprintf("%s-%s.json", r.namespace, r.key))

	entry := cacheEntry{
		Expiry:    *r.expiry,
		CreatedAt: time.Now(),
		Payload:   r.value,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
