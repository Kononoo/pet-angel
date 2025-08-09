package ai

import (
	"context"

	"github.com/go-kratos/kratos/v2/config"
)

// Config defines generic AI client configuration
type Config struct {
	Model       string  `json:"model" yaml:"model"`
	BaseURL     string  `json:"base_url" yaml:"base_url"`
	APIKey      string  `json:"api_key" yaml:"api_key"`
	MaxTokens   int     `json:"max_tokens" yaml:"max_tokens"`
	Temperature float32 `json:"temperature" yaml:"temperature"`
}

// Client is the generic AI chat interface
type Client interface {
	Chat(ctx context.Context, systemPrompt, userContent string) (string, error)
	Stream(ctx context.Context, systemPrompt, userContent string, onDelta func(string) error) (string, error)
}

var defaultClient Client

func SetClient(c Client) { defaultClient = c }
func Default() Client    { return defaultClient }

// LoadFromConfig loads "ai" section from Kratos config and constructs a client.
// If not present, it uses sensible defaults for demo.
func LoadFromConfig(c config.Config) {
	var cfg Config
	if v := c.Value("ai"); v != nil {
		_ = v.Scan(&cfg)
	}
	if cfg.Model == "" {
		cfg.Model = "deepseek-ai/DeepSeek-R1-0528-Qwen3-8B"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.siliconflow.cn/v1/"
	}
	if cfg.APIKey == "" {
		cfg.APIKey = "sk-wucfvbppymimfcrmzrtbowbnpquyudkjbjpzahlavmlhddmq"
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 16384
	}
	if cfg.Temperature == 0 {
		cfg.Temperature = 0.0
	}
	SetClient(NewClient(cfg))
}
