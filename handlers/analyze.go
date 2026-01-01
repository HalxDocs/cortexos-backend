package handlers

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"cortex/backend/config"
)

type AnalyzeRequest struct {
	Text string `json:"text"`
}

func AnalyzeThought(c *gin.Context) {
	var req AnalyzeRequest

	// -----------------------------
	// 1. Validate request body
	// -----------------------------
	if err := c.ShouldBindJSON(&req); err != nil {
		config.Log(config.ERROR, "invalid_request_body", err.Error())
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	text := strings.TrimSpace(req.Text)

	if text == "" {
		c.JSON(400, gin.H{"error": "text cannot be empty"})
		return
	}

	if len(text) > 3000 {
		c.JSON(400, gin.H{"error": "text too long (max 3000 chars)"})
		return
	}

	config.Log(config.INFO, "analysis_started", map[string]any{
		"text_length": len(text),
	})

	// -----------------------------
	// 2. Load system prompt
	// -----------------------------
	systemPrompt, err := os.ReadFile("prompts/system.txt")
	if err != nil {
		config.Log(config.ERROR, "system_prompt_missing", err.Error())
		c.JSON(500, gin.H{"error": "analysis failed"})
		return
	}

	// -----------------------------
	// 3. Call AI engine
	// -----------------------------
	response, err := config.CallOpenAI(string(systemPrompt), text)
	if err != nil {
		config.Log(config.ERROR, "openai_call_failed", err.Error())
		c.JSON(500, gin.H{"error": "analysis failed"})
		return
	}

	// -----------------------------
	// 4. Parse AI JSON (FLEXIBLE)
	// -----------------------------
	var cognition map[string]any

	decoder := json.NewDecoder(strings.NewReader(response))
	if err := decoder.Decode(&cognition); err != nil {
		config.Log(config.ERROR, "invalid_ai_json", map[string]any{
			"error": err.Error(),
		})
		c.JSON(500, gin.H{"error": "invalid AI response format"})
		return
	}

	// -----------------------------
	// 5. Normalize + respond
	// -----------------------------
	normalized := normalizeCognition(cognition)

	config.Log(config.INFO, "analysis_completed", map[string]any{
		"confidence": normalized["confidence"],
	})

	c.JSON(200, gin.H{
		"data": normalized,
		"meta": gin.H{
			"processed_at": time.Now().UTC(),
		},
	})
}
