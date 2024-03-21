package ratelimit

import (
	"sync"

	"go.uber.org/ratelimit"
)

type GlobalRateLimit struct {
	interval    int
	ratelimiter sync.Map
}

func NewGlobalRateLimit(interval int) *GlobalRateLimit {
	return &GlobalRateLimit{interval: interval, ratelimiter: sync.Map{}}
}

func (g *GlobalRateLimit) Add(key string) {
	g.ratelimiter.Store(key, ratelimit.New(g.interval))
}

func (g *GlobalRateLimit) Take(key string) {
	if rateLimit, ok := g.ratelimiter.Load(key); ok {
		rateLimit.(ratelimit.Limiter).Take()
	}
}

func (g *GlobalRateLimit) Delete(key string) {
	g.ratelimiter.Delete(key)
}

func (g *GlobalRateLimit) GetInterval() int {
	return g.interval
}
