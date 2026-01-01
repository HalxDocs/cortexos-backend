package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type LogLevel string

const (
	INFO  LogLevel = "info"
	ERROR LogLevel = "error"
)

type LogEntry struct {
	Timestamp string    `json:"timestamp"`
	Level     LogLevel `json:"level"`
	Message   string   `json:"message"`
	Context   any      `json:"context,omitempty"`
}

var logger = log.New(os.Stdout, "", 0)

func Log(level LogLevel, message string, context any) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Context:   context,
	}

	data, _ := json.Marshal(entry)
	logger.Println(string(data))
}
