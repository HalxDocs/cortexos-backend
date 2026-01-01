package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

func CallOpenAI(systemPrompt, userText string) (string, error) {
	payload, _ := json.Marshal(OpenAIRequest{
		Model: "gpt-4.1-mini",
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userText},
		},
	})

	// one retry max
	for attempt := 1; attempt <= 2; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(
			ctx,
			"POST",
			"https://api.openai.com/v1/chat/completions",
			bytes.NewBuffer(payload),
		)
		if err != nil {
			return "", err
		}

		req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
		req.Header.Set("Content-Type", "application/json")

		res, err := httpClient.Do(req)
		if err != nil {
			if attempt == 2 {
				Log(ERROR, "openai_http_failed", err.Error())
				return "", err
			}
			continue
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			Log(ERROR, "openai_read_failed", err.Error())
			return "", err
		}

		if res.StatusCode != http.StatusOK {
			Log(ERROR, "openai_non_200", map[string]any{
				"status": res.StatusCode,
				"body":   string(body),
			})
			return "", errors.New("openai request failed")
		}

		var parsed OpenAIResponse
		if err := json.Unmarshal(body, &parsed); err != nil {
			Log(ERROR, "openai_invalid_json", err.Error())
			return "", err
		}

		if len(parsed.Choices) == 0 {
			Log(ERROR, "openai_no_choices", nil)
			return "", errors.New("no response from openai")
		}

		return parsed.Choices[0].Message.Content, nil
	}

	return "", errors.New("openai request failed after retry")
}
