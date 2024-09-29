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

//go:generate mockgen -source $GOFILE -package cachetest -destination ../../../tests/mocks-gen/cachetest/inmemory_cache_mock.go

// TTLCacheClient cache methods available
type TTLCacheClient interface {
	Get(key string) (interface{}, bool)
	SetWithTTL(key string, data interface{}, ttl time.Duration)
	Remove(key string) bool
}

// Caching method available for cache
type Caching interface {
	Get(key string) (interface{}, bool)
	Set(key, value string, ttl time.Duration) bool
	Delete(key string)
}
