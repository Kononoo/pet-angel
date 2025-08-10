package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	return &sfClient{http: &http.Client{Timeout: 120 * time.Second}, cfg: cfg}
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

	// 使用独立的上下文，避免HTTP请求的超时限制
	reqCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(reqCtx, http.MethodPost, strings.TrimRight(c.cfg.BaseURL, "/")+"/chat/completions", bytes.NewReader(buf))
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// 打印请求信息用于调试
	fmt.Printf("AI Chat Request: URL=%s, Model=%s, APIKey=%s...\n", req.URL.String(), c.cfg.Model, c.cfg.APIKey[:10])

	resp, err := c.http.Do(req)
	if err != nil {
		fmt.Printf("AI Chat HTTP Error: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("AI Chat HTTP Error: status=%d, body=%s\n", resp.StatusCode, string(bodyBytes))
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var out sfResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		fmt.Printf("AI Chat JSON Decode Error: %v\n", err)
		return "", err
	}
	if len(out.Choices) == 0 {
		fmt.Printf("AI Chat Empty Choices\n")
		return "", errors.New("empty choices")
	}

	content := out.Choices[0].Message.Content
	fmt.Printf("AI Chat Response: length=%d, content=%s\n", len(content), content)
	return content, nil
}

// Stream streaming; onDelta will be called with each text delta
func (c *sfClient) Stream(ctx context.Context, systemPrompt, userContent string, onDelta func(string) error) (string, error) {
	body := &sfReq{Model: c.cfg.Model, Messages: []sfMsg{{Role: "system", Content: systemPrompt}, {Role: "user", Content: userContent}}, Stream: true, MaxTokens: c.cfg.MaxTokens}
	buf, _ := json.Marshal(body)

	// 使用独立的上下文，避免HTTP请求的超时限制
	reqCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(reqCtx, http.MethodPost, strings.TrimRight(c.cfg.BaseURL, "/")+"/chat/completions", bytes.NewReader(buf))
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// 打印请求信息用于调试
	fmt.Printf("AI Stream Request: URL=%s, Model=%s, APIKey=%s...\n", req.URL.String(), c.cfg.Model, c.cfg.APIKey[:10])

	resp, err := c.http.Do(req)
	if err != nil {
		fmt.Printf("AI Stream HTTP Error: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("AI Stream HTTP Error: status=%d, body=%s\n", resp.StatusCode, string(bodyBytes))
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	reader := bufio.NewReader(resp.Body)
	var full strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Printf("AI Stream Read Error: %v\n", err)
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

	fmt.Printf("AI Stream Response: length=%d, content=%s\n", full.Len(), full.String())
	return full.String(), nil
}
