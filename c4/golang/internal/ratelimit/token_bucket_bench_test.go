package ratelimit

import (
	"testing"
	"time"
)

func BenchmarkTokenBucketMutex(b *testing.B) {
	t := NewTokenBucketMutex(250, 10, 5*time.Microsecond)
	for i := 0; i < b.N; i++ {
		t.Allow("some_key")
	}
}

func BenchmarkTokenBucketCAS(b *testing.B) {
	t := NewTokenBucketCas(250, 10, 5*time.Microsecond)
	for i := 0; i < b.N; i++ {
		t.Allow("some_key")
	}
}

// BenchmarkTokenBucket_MultiKey tests performance when requests are spread across many keys
// This simulates realistic API usage where different users/IPs are rate limited independently
// Expected: CAS might perform better due to less lock contention across different buckets
// TODO: Implement by cycling through multiple keys (e.g., key_0, key_1, ..., key_99)
// func BenchmarkTokenBucket_MultiKey(b *testing.B) {}

// BenchmarkTokenBucket_Parallel tests performance under concurrent access from multiple goroutines
// This simulates real-world concurrent requests hitting the rate limiter
// Expected: Should reveal if either implementation has better parallelism characteristics
// TODO: Implement using b.RunParallel() and pb.Next() pattern
// func BenchmarkTokenBucket_Parallel(b *testing.B) {}

// BenchmarkTokenBucket_HighContention tests performance when most requests are denied (tokens exhausted)
// This simulates what happens during an actual rate limiting event (user hitting the limit)
// Expected: Should be faster than allow case since we just check tokens <= 0 and return
// TODO: Implement with low capacity (e.g., 1) so most calls return false
// func BenchmarkTokenBucket_HighContention(b *testing.B) {}
