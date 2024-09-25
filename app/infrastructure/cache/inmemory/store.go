package inmemory

import (
	"time"

	"github.com/ReneKroon/ttlcache"

	"github.com/shubhamgptln/hotels-data-merge/app/infrastructure/cache"
)

// Service ...
type Service struct {
	client cache.TTLCacheClient
}

// New returns new cache interface
func NewInMemoryCache() *Service {
	c := ttlcache.NewCache()
	return &Service{
		client: c,
	}
}

// Get to get key in cache
func (c Service) Get(key string) (interface{}, bool) {
	return c.client.Get(key)
}

// Set to set cache with ttl
func (c Service) Set(key, value string, ttl time.Duration) bool {
	c.client.SetWithTTL(key, value, ttl)
	return true
}

// Delete to delete key in cache
func (c Service) Delete(key string) {
	c.client.Remove(key)
}
