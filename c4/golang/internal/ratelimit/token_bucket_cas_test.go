package ratelimit

import (
	"testing"
	"time"
)

func TestTokenBucketCas(t *testing.T) {
	capacity := 3
	tb := NewTokenBucketCas(int32(capacity), 1, time.Hour)
	for i := range capacity {
		if !tb.Allow("me") {
			t.Fatalf("attempt %d should be allowed", i)
		}
	}
	if tb.Allow("me") {
		t.Fatalf("attempt %d should be disallowed", capacity+1)
	}
}

func TestTokenBucketCas_Refill(t *testing.T) {
	// 2 tokens, refill 1 token every 100ms
	tb := NewTokenBucketCas(2, 1, 10*time.Millisecond)

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
}

func TestTokenBucketCas_Multikey(t *testing.T) {
	capacity := 3
	tb := NewTokenBucketCas(int32(capacity), 1, time.Hour)
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

}
