package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type SimpleMessage struct {
	Text string `json:"text"`
}

type MessageClient struct {
	webhookURL string
	platform   string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewMessageClient(webhookURL, platform string, logger *zap.Logger) *MessageClient {
	return &MessageClient{
		webhookURL: webhookURL,
		platform:   platform,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (mc *MessageClient) PostMessage(msg *SimpleMessage) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", mc.webhookURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := mc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%s returned status %d: %s", mc.platform, resp.StatusCode, string(body))
	}

	return nil
}

func BuildMessage(alert *GraylogAlert, cfg *Config) *SimpleMessage {
	severity := alert.GetSeverityName()
	message := alert.GetDisplayMessage()
	timestamp := alert.GetTimestamp()

	if len(message) > 500 {
		message = message[:500] + "..."
	}

	text := fmt.Sprintf("[%s] %s\nTime: %s", severity, message, timestamp.Format(time.RFC3339))
	
	if alert.Source != "" {
		text += fmt.Sprintf("\nSource: %s", alert.Source)
	}
	
	if alert.EventDefinitionID != "" {
		text += fmt.Sprintf("\nEvent ID: %s", alert.EventDefinitionID)
	}

	msg := &SimpleMessage{
		Text: text,
	}

	return msg
}

func getSeverityColor(severity string) string {
	switch strings.ToLower(severity) {
	case "0", "1":
		return "#D62828"
	case "2":
		return "#F77F00"
	case "3":
		return "#FFB703"
	case "4", "5":
		return "#219EBC"
	case "6":
		return "#023047"
	default:
		return "#999999"
	}
}