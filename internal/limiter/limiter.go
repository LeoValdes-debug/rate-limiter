package limiter

import (
	"sync"
	"time"
)

type bucket struct {
	tokens    float64
	lastRefil time.Time
	mu        sync.Mutex
}

func newBucket(capacity float64) *bucket {
	return &bucket{
		tokens:    capacity,
		lastRefil: time.Now(),
	}
}

func (b *bucket) allow(capacity, refillRate float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefil).Seconds()
	b.tokens += elapsed * refillRate
	if b.tokens > capacity {
		b.tokens = capacity
	}
	b.lastRefil = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// Limiter ограничивает запросы по IP.
// Использует алгоритм token bucket: каждый IP получает N токенов,
// токены пополняются со временем.
type Limiter struct {
	mu         sync.Mutex
	buckets    map[string]*bucket
	capacity   float64
	refillRate float64 // токенов в секунду
}

func New(capacity float64, refillRate float64) *Limiter {
	l := &Limiter{
		buckets:    make(map[string]*bucket),
		capacity:   capacity,
		refillRate: refillRate,
	}
	go l.cleanup()
	return l
}

func (l *Limiter) Allow(ip string) bool {
	l.mu.Lock()
	b, ok := l.buckets[ip]
	if !ok {
		b = newBucket(l.capacity)
		l.buckets[ip] = b
	}
	l.mu.Unlock()

	return b.allow(l.capacity, l.refillRate)
}

// cleanup удаляет старые бакеты раз в минуту
func (l *Limiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		l.mu.Lock()
		for ip, b := range l.buckets {
			b.mu.Lock()
			idle := time.Since(b.lastRefil) > 5*time.Minute
			b.mu.Unlock()
			if idle {
				delete(l.buckets, ip)
			}
		}
		l.mu.Unlock()
	}
}
