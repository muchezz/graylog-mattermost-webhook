package main

import (
	"io"
	"net/http"

	"go.uber.org/zap"
)

type GraylogHandler struct {
	config          *Config
	logger          *zap.Logger
	messageClient   *MessageClient
}

func NewGraylogHandler(cfg *Config, logger *zap.Logger) *GraylogHandler {
	return &GraylogHandler{
		config:        cfg,
		logger:        logger,
		messageClient: NewMessageClient(cfg.Destination.WebhookURL, cfg.Destination.Platform, logger),
	}
}

func (gh *GraylogHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		gh.logger.Warn("Invalid request method", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		gh.logger.Error("Failed to read request body", zap.Error(err))
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Parse the alert
	alert, err := ParseGraylogAlert(body)
	if err != nil {
		gh.logger.Error("Failed to parse Graylog alert", zap.Error(err))
		http.Error(w, "Failed to parse alert", http.StatusBadRequest)
		return
	}

	gh.logger.Info("Received alert",
		zap.String("event_id", alert.EventDefinitionID),
		zap.String("severity", alert.GetSeverityName()),
		zap.String("message", truncate(alert.GetDisplayMessage(), 100)),
	)

	// Build message
	message := BuildMessage(alert, gh.config)

	// Post to Slack or Mattermost
	if err := gh.messageClient.PostMessage(message); err != nil {
		gh.logger.Error("Failed to post message",
			zap.Error(err),
			zap.String("platform", gh.config.Destination.Platform),
			zap.String("channel", message.Channel),
		)
		http.Error(w, "Failed to post message", http.StatusInternalServerError)
		return
	}

	gh.logger.Info("Alert posted successfully",
		zap.String("event_id", alert.EventDefinitionID),
		zap.String("platform", gh.config.Destination.Platform),
		zap.String("channel", message.Channel),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
