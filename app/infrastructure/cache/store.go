package cache

import (
	"time"
)

const (
	// ErrUnableToStoreCache unable to store cache
	ErrUnableToStoreCache = "unable to store cache"
	// ErrCacheNotFound object is not found
	ErrCacheNotFound = "object is not found"
)

// TTLCacheClient cache methods available
type TTLCacheClient interface {
	Get(key string) (interface{}, bool)
	SetWithTTL(key string, data interface{}, ttl time.Duration)
	Remove(key string) bool
}

// Methods method available for cache
type Methods interface {
	Get(key string) (interface{}, bool)
	Set(key, value string, ttl time.Duration) bool
	Delete(key string)
}
