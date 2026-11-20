package crawler

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// URLFrontierCache wraps a URL frontier and adds caching functionality
type URLFrontierCache struct {
	cache      *cache.Cache
	urlFrontier URLFrontier
}

// NewURLFrontierCache returns a new URLFrontierCache with the given frontier and cache options
func NewURLFrontierCache(urlFrontier URLFrontier, expiration, cleanupInterval time.Duration) *URLFrontierCache {
	return &URLFrontierCache{
		cache:      cache.New(expiration, cleanupInterval),
		urlFrontier: urlFrontier,
	}
}

// ReadFromCache fetches a URL from the cache or the underlying frontier
func (f *URLFrontierCache) ReadFromCache() (string, bool) {
	url, ok := f.cache.Get("url")
	if ok {
		return url.(string), true
	}
	url, ok = f.urlFrontier.Read()
	if ok {
		f.cache.Set("url", url, cache.DefaultExpiration)
		return url, true
	}
	return "", false
}

// Has fetches a URL from the cache or the underlying frontier and checks if it has more URLs
func (f *URLFrontierCache) Has() bool {
	url, ok := f.ReadFromCache()
	return ok && f.urlFrontier.HasMore(url)
}

// Discard discards a URL from the cache and the underlying frontier
func (f *URLFrontierCache) Discard() {
	f.cache.Delete("url")
	f.urlFrontier.Discard()
}
