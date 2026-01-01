package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().UTC(),
	})
}

func Ready(c *gin.Context) {
	// later we’ll check DB, cache, OpenAI, etc.
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
