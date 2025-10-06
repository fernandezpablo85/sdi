package ratelimit

type RateLimiter interface {
	Allow(key string) bool
}
