package rate_limiter

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrLimitExceeded = errors.New("rate limit exceeded")

type RequestUserInfo struct {
	Count       int
	RequestedAt time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]RequestUserInfo
}

func New() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string]RequestUserInfo),
	}
}
func (rl *RateLimiter) Allow(key string) (bool, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	info, exists := rl.requests[key]

	if !exists || now.Sub(info.RequestedAt) >= 1*time.Minute {
		rl.requests[key] = RequestUserInfo{
			Count:       1,
			RequestedAt: now,
		}
		return true, nil
	}

	if info.Count >= 5 {
		return false, ErrLimitExceeded
	}
	info.Count++
	rl.requests[key] = info
	return true, nil
}
func (rl *RateLimiter) WorkerClear(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for key, info := range rl.requests {
				if now.Sub(info.RequestedAt) >= 1*time.Minute {
					delete(rl.requests, key)
				}
			}
			rl.mu.Unlock()
		}
	}
}
