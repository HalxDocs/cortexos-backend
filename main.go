package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"cortex/backend/config"
	"cortex/backend/handlers"
	"cortex/backend/invites"
)

func init() {
	// ----------------------------------
	// Invite codes (DEV / PRIVATE BETA)
	// ----------------------------------
	invites.CreateInvite("alpha-1", 7*24*time.Hour)
	invites.CreateInvite("alpha-2", 7*24*time.Hour)
	invites.CreateInvite("alpha-3", 7*24*time.Hour)
}

func main() {
	r := gin.New()

	// ----------------------------------
	// 0. CORS (MUST BE FIRST)
	// ----------------------------------
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost")
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	// ----------------------------------
	// 1. Rate Limiter
	// ----------------------------------
	limiter := config.NewRateLimiter(10, time.Minute)
	r.Use(limiter.Middleware())

	// ----------------------------------
	// 2. Structured Logging
	// ----------------------------------
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		config.Log(config.INFO, "http_request", map[string]any{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency_ms": time.Since(start).Milliseconds(),
			"ip":         c.ClientIP(),
		})
	})

	// ----------------------------------
	// 3. Panic Recovery
	// ----------------------------------
	r.Use(gin.Recovery())

	// ----------------------------------
	// 4. Global Routes
	// ----------------------------------
	r.GET("/health", handlers.Health)
	r.GET("/ready", handlers.Ready)

	// ----------------------------------
	// 5. Versioned API Routes
	// ----------------------------------
	v1 := r.Group("/v1")
	{
		// Invite redemption
		v1.POST("/invite/:code", handlers.RedeemInvite)

		// Thought analysis
		v1.OPTIONS("/analyze-thought", func(c *gin.Context) {
			c.Status(204)
		})
		v1.POST("/analyze-thought", handlers.AnalyzeThought)
	}

	// ----------------------------------
	// 6. Start Server
	// ----------------------------------
	r.Run(":8080")
}
