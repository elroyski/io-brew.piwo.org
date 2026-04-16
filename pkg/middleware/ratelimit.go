package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ipRecord struct {
	count       int
	windowStart time.Time
}

// RateLimiter ogranicza liczbę żądań per IP w danym oknie czasowym
type RateLimiter struct {
	mu      sync.Mutex
	records map[string]*ipRecord
	max     int
	window  time.Duration
}

// NewRateLimiter tworzy nowy rate limiter; max to liczba dozwolonych żądań w oknie window
func NewRateLimiter(max int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		records: make(map[string]*ipRecord),
		max:     max,
		window:  window,
	}
	// Goroutine czyszcząca stare wpisy
	go func() {
		for {
			time.Sleep(window)
			rl.mu.Lock()
			now := time.Now()
			for ip, rec := range rl.records {
				if now.Sub(rec.windowStart) > window {
					delete(rl.records, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
	return rl
}

// Allow sprawdza czy IP ma jeszcze wolne "sloty" w bieżącym oknie
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	rec, exists := rl.records[ip]
	if !exists || now.Sub(rec.windowStart) > rl.window {
		rl.records[ip] = &ipRecord{count: 1, windowStart: now}
		return true
	}

	if rec.count >= rl.max {
		return false
	}
	rec.count++
	return true
}

// RegisterMiddleware zwraca middleware Gin dla endpointu rejestracji
func (rl *RateLimiter) RegisterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.Allow(ip) {
			c.HTML(http.StatusTooManyRequests, "register.html", gin.H{
				"error": "Zbyt wiele prób rejestracji z tego adresu IP. Spróbuj ponownie za godzinę.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
