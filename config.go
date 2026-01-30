package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Mattermost MattermostConfig `yaml:"mattermost"`
}

type ServerConfig struct {
	ListenAddr string `yaml:"listen_addr"`
	LogLevel   string `yaml:"log_level"`
}

type MattermostConfig struct {
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
		Mattermost: MattermostConfig{
			Username:  "Graylog",
			IconEmoji: ":clipboard:",
		},
	}

	// Try to load from config file
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "/etc/graylog-mattermost-webhook/config.yaml"
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
	if webhookURL := os.Getenv("MATTERMOST_WEBHOOK_URL"); webhookURL != "" {
		cfg.Mattermost.WebhookURL = webhookURL
	}
	if username := os.Getenv("MATTERMOST_USERNAME"); username != "" {
		cfg.Mattermost.Username = username
	}
	if channel := os.Getenv("MATTERMOST_CHANNEL"); channel != "" {
		cfg.Mattermost.Channel = channel
	}

	// Validate required configuration
	if cfg.Mattermost.WebhookURL == "" {
		return nil, fmt.Errorf("MATTERMOST_WEBHOOK_URL environment variable or config is required")
	}

	return cfg, nil
}
