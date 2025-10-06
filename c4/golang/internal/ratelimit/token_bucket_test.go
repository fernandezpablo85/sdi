package ratelimit

import (
	"testing"
	"time"
)

var implementations = []struct {
	name    string
	factory func(capacity, refillAmount int32, refillPeriod time.Duration) RateLimiter
}{
	{"CAS", func(c, r int32, p time.Duration) RateLimiter { return NewTokenBucketCas(c, r, p) }},
	{"Mutex", func(c, r int32, p time.Duration) RateLimiter { return NewTokenBucketMutex(c, r, p) }},
}

func TestTokenBucket(t *testing.T) {
	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			capacity := 3
			tb := impl.factory(int32(capacity), 1, time.Hour)
			for i := range capacity {
				if !tb.Allow("me") {
					t.Fatalf("attempt %d should be allowed", i)
				}
			}
			if tb.Allow("me") {
				t.Fatalf("attempt %d should be disallowed", capacity+1)
			}
		})
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			// 2 tokens, refill 1 token every 10ms
			tb := impl.factory(2, 1, 10*time.Millisecond)

			// Consume all tokens
			tb.Allow("test")
			tb.Allow("test")

			// Should be denied now
			if tb.Allow("test") {
				t.Error("Should be denied when empty")
			}

			// Wait for refill
			time.Sleep(50 * time.Millisecond)

			// Should allow 1 request after refill
			if !tb.Allow("test") {
				t.Error("Should allow after refill")
			}
		})
	}
}

func TestTokenBucket_Multikey(t *testing.T) {
	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			capacity := 3
			tb := impl.factory(int32(capacity), 1, time.Hour)
			for i := range capacity {
				if !tb.Allow("me") {
					t.Fatalf("attempt %d should be allowed", i)
				}
			}
			for i := range capacity {
				if !tb.Allow("other") {
					t.Fatalf("attempt %d should be allowed", i)
				}
			}
			if tb.Allow("me") {
				t.Fatalf("key 'me' should be exhausted at attempt %d", capacity+1)
			}
			if tb.Allow("other") {
				t.Fatalf("key 'other' should be exhausted at attempt %d", capacity+1)
			}
		})
	}
}
