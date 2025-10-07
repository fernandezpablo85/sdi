package ratelimit

import (
	"sync"
	"sync/atomic"
	"time"
)

type bucket struct {
	tokens atomic.Int32
}

type TokenBucketCas struct {
	buckets      sync.Map
	capacity     int32
	refillAmount int32
	refillPeriod time.Duration
}

func (t *TokenBucketCas) Allow(key string) bool {
	if val, ok := t.buckets.Load(key); ok {
		b := val.(*bucket)
		return b.allow()
	}
	newBucket := t.newBucket()
	actual, _ := t.buckets.LoadOrStore(key, newBucket)
	b := actual.(*bucket)
	return b.allow()
}

func (t *TokenBucketCas) newBucket() *bucket {
	b := bucket{}
	b.tokens.Store(t.capacity)
	return &b
}

func (b *bucket) allow() bool {
	for {
		current := b.tokens.Load()
		if current <= 0 {
			return false
		}

		ok := b.tokens.CompareAndSwap(current, current-1)
		if ok {
			return true
		}
	}
}

func NewTokenBucketCas(capacity, refillAmount int32, refillPeriod time.Duration) *TokenBucketCas {
	t := &TokenBucketCas{capacity: capacity, refillAmount: refillAmount, refillPeriod: refillPeriod}
	go func() {
		for range time.Tick(t.refillPeriod) {
			t.buckets.Range(func(key, value any) bool {
				for {
					b := value.(*bucket)
					current := b.tokens.Load()
					updated := min(t.capacity, current+t.refillAmount)
					ok := b.tokens.CompareAndSwap(current, updated)
					if ok {
						break
					}
				}
				return true
			})
		}
	}()
	return t
}
