package ratelimiter

import (
	"errors"
	"sync"
	"time"
)

type RateLimiter struct{
 mu sync.Mutex
 requests map[int]RequestUserInfo
}

type RequestUserInfo struct{
 Counter int
 RequestedAt time.Time
}

func New() *RateLimiter{
 return &RateLimiter{
  mu: sync.Mutex{},
  requests: make(map[int]RequestUserInfo),
 }
}

func(rl *RateLimiter) Allow(id int) (bool, error){
	rl.mu.Lock()
	defer rl.mu.Unlock()
	info, exist:=rl.requests[id]
	now:=time.Now()
	if !exist{
		rl.requests[id] = RequestUserInfo{
			Counter: 1,
			RequestedAt: now,
		}
		return true, nil
	}

	if now.Sub(info.RequestedAt)>=time.Minute*1{
		rl.requests[id]=RequestUserInfo{
			Counter: 1,
			RequestedAt: now,
		}
		return true, nil
	}


	if info.Counter>=5{
		return false, errors.New("limit has been exceeded")
	}
	info.Counter++
	rl.requests[id]=info
	return true, nil
}

func(rl *RateLimiter)WorkerClear(){
	tick:=time.NewTicker(5*time.Second)
	for range tick.C{
		rl.mu.Lock()
		now:=time.Now()	
		for id, info:=range rl.requests{
			if now.Sub(info.RequestedAt)>=time.Minute*1{
				delete(rl.requests, id)
			}
		}
		rl.mu.Unlock()
	}
}

