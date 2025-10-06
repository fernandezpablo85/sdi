package ratelimit

import (
	"sync"
	"time"
)

type lockBucket struct {
	tokens int32
}

type TokenBucketMutex struct {
	buckets      map[string]*lockBucket
	lock         sync.RWMutex
	capacity     int32
	refillAmount int32
	refillPeriod time.Duration
}

func NewTokenBucketMutex(capacity, refillAmount int32, refillPeriod time.Duration) *TokenBucketMutex {
	return &TokenBucketMutex{
		capacity:     capacity,
		refillAmount: refillAmount,
		refillPeriod: refillPeriod,
		buckets:      make(map[string]*lockBucket),
	}
}

func (t *TokenBucketMutex) newLockBucket() *lockBucket {
	b := &lockBucket{tokens: t.capacity}
	go func() {
		for range time.Tick(t.refillPeriod) {
			t.lock.Lock()
			b.tokens = min(t.capacity, b.tokens+t.refillAmount)
			t.lock.Unlock()
		}
	}()
	return b
}

func (t *TokenBucketMutex) Allow(key string) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	bucket, present := t.buckets[key]
	if !present {
		bucket = t.newLockBucket()
		t.buckets[key] = bucket
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	} else {
		return false
	}
}
