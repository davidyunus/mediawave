package tokenbucket

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int
	tokens     int
	refillRate int // tokens per second
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) AllowRequest() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	refill := int(elapsed * float64(tb.refillRate))

	if refill > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+refill)
		tb.lastRefill = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Simulasi pemanggilan
func main() {
	bucket := NewTokenBucket(5, 2) // capacity 5 token, refill 2 token per detik

	for i := 0; i < 10; i++ {
		allowed := bucket.AllowRequest()
		fmt.Printf("Request #%d allowed? %v\n", i+1, allowed)
		time.Sleep(300 * time.Millisecond)
	}
}
