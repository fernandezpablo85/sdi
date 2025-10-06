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
	val, _ := t.buckets.LoadOrStore(key, t.newBucket())
	b := val.(*bucket)
	return b.allow()
}

func (t *TokenBucketCas) newBucket() *bucket {
	b := bucket{}
	b.tokens.Store(t.capacity)
	go func() {
		for range time.Tick(t.refillPeriod) {
			for {
				current := b.tokens.Load()
				updated := min(t.capacity, current+t.refillAmount)
				ok := b.tokens.CompareAndSwap(current, updated)
				if ok {
					break
				}
			}
		}
	}()
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
	return &TokenBucketCas{capacity: capacity, refillAmount: refillAmount, refillPeriod: refillPeriod}
}
