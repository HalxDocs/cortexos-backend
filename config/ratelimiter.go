package config

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type client struct {
	requests int
	lastSeen time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
	}

	// cleanup goroutine
	go func() {
		for {
			time.Sleep(window)
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, c := range rl.clients {
		if now.Sub(c.lastSeen) > rl.window {
			delete(rl.clients, ip)
		}
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ✅ VERY IMPORTANT:
		// Allow CORS preflight requests to pass through
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		ip := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		cData, exists := rl.clients[ip]
		if !exists {
			rl.clients[ip] = &client{
				requests: 1,
				lastSeen: time.Now(),
			}
			c.Next()
			return
		}

		// reset window
		if time.Since(cData.lastSeen) > rl.window {
			cData.requests = 1
			cData.lastSeen = time.Now()
			c.Next()
			return
		}

		if cData.requests >= rl.limit {
			Log(ERROR, "rate_limit_exceeded", map[string]any{
				"ip": ip,
			})

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		cData.requests++
		cData.lastSeen = time.Now()
		c.Next()
	}
}
