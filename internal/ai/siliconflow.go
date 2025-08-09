package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// SiliconFlow Chat Completions minimal client (non-stream + stream)
// Docs: https://docs.siliconflow.cn/cn/api-reference/chat-completions/chat-completions

type sfClient struct {
	http *http.Client
	cfg  Config
}

func NewClient(cfg Config) Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.siliconflow.cn/v1/"
	}
	if cfg.Model == "" {
		cfg.Model = "deepseek-ai/DeepSeek-R1-0528-Qwen3-8B"
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 16384
	}
	return &sfClient{http: &http.Client{Timeout: 60 * time.Second}, cfg: cfg}
}

type sfMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type sfReq struct {
	Model     string  `json:"model"`
	Messages  []sfMsg `json:"messages"`
	Stream    bool    `json:"stream,omitempty"`
	MaxTokens int     `json:"max_tokens,omitempty"`
}
type sfResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Chat non-streaming
func (c *sfClient) Chat(ctx context.Context, systemPrompt, userContent string) (string, error) {
	body := &sfReq{Model: c.cfg.Model, Messages: []sfMsg{{Role: "system", Content: systemPrompt}, {Role: "user", Content: userContent}}, MaxTokens: c.cfg.MaxTokens}
	buf, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.cfg.BaseURL, "/")+"/chat/completions", bytes.NewReader(buf))
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out sfResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("empty choices")
	}
	return out.Choices[0].Message.Content, nil
}

// Stream streaming; onDelta will be called with each text delta
func (c *sfClient) Stream(ctx context.Context, systemPrompt, userContent string, onDelta func(string) error) (string, error) {
	body := &sfReq{Model: c.cfg.Model, Messages: []sfMsg{{Role: "system", Content: systemPrompt}, {Role: "user", Content: userContent}}, Stream: true, MaxTokens: c.cfg.MaxTokens}
	buf, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.cfg.BaseURL, "/")+"/chat/completions", bytes.NewReader(buf))
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	reader := bufio.NewReader(resp.Body)
	var full strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return full.String(), err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			break
		}
		var obj map[string]any
		if json.Unmarshal([]byte(payload), &obj) == nil {
			// OpenAI-style: choices[0].delta.content
			if ch, ok := obj["choices"].([]any); ok && len(ch) > 0 {
				if m, ok := ch[0].(map[string]any); ok {
					if delta, ok := m["delta"].(map[string]any); ok {
						if txt, ok := delta["content"].(string); ok && txt != "" {
							full.WriteString(txt)
							if onDelta != nil {
								_ = onDelta(txt)
							}
						}
					}
				}
			}
		}
	}
	return full.String(), nil
}
