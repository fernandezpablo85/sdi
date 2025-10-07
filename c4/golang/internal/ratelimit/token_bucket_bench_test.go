package ratelimit

import (
	"fmt"
	"math/rand/v2"
	"runtime"
	"testing"
	"time"
)

func BenchmarkTokenBucketMutex_SingleKey(b *testing.B) {
	runtime.GC()
	b.Run("Mutex", func(b *testing.B) {
		t := NewTokenBucketMutex(250, 10, 5*time.Microsecond)
		for i := 0; i < b.N; i++ {
			t.Allow("some_key")
		}
	})

	b.Run("Cas", func(b *testing.B) {
		t := NewTokenBucketCas(250, 10, 5*time.Microsecond)
		for i := 0; i < b.N; i++ {
			t.Allow("some_key")
		}
	})

}

func BenchmarkTokenBucket_MultiKey(b *testing.B) {
	runtime.GC()
	n := 100
	keys := make([]string, n)
	for i := range n {
		keys[i] = fmt.Sprintf("key_%d", i)
	}

	b.Run("Mutex", func(b *testing.B) {
		t := NewTokenBucketMutex(250, 10, 5*time.Microsecond)
		// Warmup: pre-create all buckets
		for i := 0; i < n; i++ {
			t.Allow(keys[i])
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ix := rand.Int() % n
			t.Allow(keys[ix])
		}
	})

	b.Run("Cas", func(b *testing.B) {
		t := NewTokenBucketCas(250, 10, 5*time.Microsecond)
		// Warmup: pre-create all buckets
		for i := 0; i < n; i++ {
			t.Allow(keys[i])
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ix := rand.Int() % n
			t.Allow(keys[ix])
		}
	})

}

func BenchmarkTokenBucket_Parallel(b *testing.B) {
	runtime.GC()
	b.Run("Mutex", func(b *testing.B) {
		t := NewTokenBucketMutex(250, 10, 5*time.Microsecond)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				t.Allow("some_key")
			}
		})
	})

	b.Run("Cas", func(b *testing.B) {
		t := NewTokenBucketCas(250, 10, 5*time.Microsecond)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				t.Allow("some_key")
			}
		})
	})
}

func BenchmarkTokenBucket_HighContention(b *testing.B) {
	runtime.GC()
	b.Run("Mutex", func(b *testing.B) {
		t := NewTokenBucketMutex(5, 10, 50*time.Microsecond)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				t.Allow("some_key")
			}
		})
	})

	b.Run("Cas", func(b *testing.B) {
		t := NewTokenBucketCas(5, 10, 50*time.Microsecond)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				t.Allow("some_key")
			}
		})
	})
}
