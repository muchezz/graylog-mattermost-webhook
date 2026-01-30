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
	Text       string        `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Fallback string `json:"fallback"`
	Color    string `json:"color"`
	Title    string `json:"title,omitempty"`
	Text     string `json:"text,omitempty"`
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
	color := getSeverityColor(alert.GetSeverity())
	message := alert.GetDisplayMessage()

	if len(message) > 500 {
		message = message[:500] + "..."
	}

	severity := alert.GetSeverityName()
	timestamp := alert.GetTimestamp()

	text := fmt.Sprintf("**[%s]** %s\n**Time:** %s", severity, message, timestamp.Format(time.RFC3339))
	
	if alert.Source != "" {
		text += fmt.Sprintf("\n**Source:** %s", alert.Source)
	}
	
	if alert.EventDefinitionID != "" {
		text += fmt.Sprintf("\n**Event ID:** %s", alert.EventDefinitionID)
	}

	attachment := Attachment{
		Fallback: fmt.Sprintf("[%s] %s", severity, message),
		Color:    color,
		Title:    message,
		Text:     text,
	}

	msg := &SimpleMessage{
		Text:        fmt.Sprintf("[%s] %s", severity, message),
		Attachments: []Attachment{attachment},
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