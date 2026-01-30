package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Destination DestinationConfig `yaml:"destination"`
}

type ServerConfig struct {
	ListenAddr string `yaml:"listen_addr"`
	LogLevel   string `yaml:"log_level"`
}

type DestinationConfig struct {
	// Platform: "slack" or "mattermost"
	Platform     string            `yaml:"platform"`
	WebhookURL   string            `yaml:"webhook_url"`
	Username     string            `yaml:"username"`
	IconEmoji    string            `yaml:"icon_emoji"`
	Channel      string            `yaml:"channel"`
	Destinations map[string]string `yaml:"destinations"` // severity -> channel mapping
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			ListenAddr: "0.0.0.0:8080",
			LogLevel:   "info",
		},
		Destination: DestinationConfig{
			Platform:  "mattermost",
			Username:  "Graylog",
			IconEmoji: ":clipboard:",
		},
	}

	// Try to load from config file
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "/etc/graylog-webhook/config.yaml"
	}

	if data, err := os.ReadFile(configPath); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables
	if addr := os.Getenv("LISTEN_ADDR"); addr != "" {
		cfg.Server.ListenAddr = addr
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.Server.LogLevel = logLevel
	}
	if platform := os.Getenv("PLATFORM"); platform != "" {
		cfg.Destination.Platform = platform
	}
	if webhookURL := os.Getenv("WEBHOOK_URL"); webhookURL != "" {
		cfg.Destination.WebhookURL = webhookURL
	}
	if username := os.Getenv("USERNAME"); username != "" {
		cfg.Destination.Username = username
	}
	if channel := os.Getenv("CHANNEL"); channel != "" {
		cfg.Destination.Channel = channel
	}

	// Validate required configuration
	if cfg.Destination.WebhookURL == "" {
		return nil, fmt.Errorf("WEBHOOK_URL environment variable or config is required")
	}
	if cfg.Destination.Platform != "slack" && cfg.Destination.Platform != "mattermost" {
		return nil, fmt.Errorf("platform must be 'slack' or 'mattermost', got: %s", cfg.Destination.Platform)
	}

	return cfg, nil
}
