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

// SimpleMessage for basic Mattermost/Slack compatibility
type SimpleMessage struct {
	Text        string                 `json:"text"`
	Channel     string                 `json:"channel,omitempty"`
	Username    string                 `json:"username,omitempty"`
	IconEmoji   string                 `json:"icon_emoji,omitempty"`
	Attachments []SimpleAttachment      `json:"attachments,omitempty"`
	Props       map[string]interface{} `json:"props,omitempty"`
}

// SimpleAttachment for message formatting
type SimpleAttachment struct {
	Fallback string                 `json:"fallback"`
	Color    string                 `json:"color"`
	Title    string                 `json:"title,omitempty"`
	Text     string                 `json:"text,omitempty"`
	Fields   []SimpleField           `json:"fields,omitempty"`
	Ts       int64                  `json:"ts,omitempty"`
}

// SimpleField for attachment fields
type SimpleField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// MessageClient handles posting to Slack or Mattermost
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

// PostMessage posts a message to Slack or Mattermost
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

// BuildMessage creates a formatted message for Slack or Mattermost
func BuildMessage(alert *GraylogAlert, cfg *Config) *SimpleMessage {
	color := getSeverityColor(alert.GetSeverity())
	timestamp := alert.GetTimestamp()
	message := alert.GetDisplayMessage()

	// Truncate message if too long
	if len(message) > 500 {
		message = message[:500] + "..."
	}

	// Build fields
	fields := []SimpleField{
		{
			Title: "Severity",
			Value: alert.GetSeverityName(),
			Short: true,
		},
		{
			Title: "Time",
			Value: timestamp.Format(time.RFC3339),
			Short: true,
		},
	}

	if alert.Source != "" {
		fields = append(fields, SimpleField{
			Title: "Source",
			Value: alert.Source,
			Short: true,
		})
	}

	if alert.EventDefinitionID != "" {
		fields = append(fields, SimpleField{
			Title: "Event ID",
			Value: alert.EventDefinitionID,
			Short: true,
		})
	}

	// Determine channel
	channel := cfg.Destination.Channel
	if cfg.Destination.Destinations != nil {
		if dest, exists := cfg.Destination.Destinations[alert.GetSeverity()]; exists {
			channel = dest
		}
	}

	// Build attachment
	attachment := SimpleAttachment{
		Fallback: fmt.Sprintf("[%s] %s", alert.GetSeverityName(), message),
		Color:    color,
		Title:    message,
		Fields:   fields,
		Ts:       timestamp.Unix(),
	}

	// Build message
	msg := &SimpleMessage{
		Channel:     channel,
		Username:    cfg.Destination.Username,
		IconEmoji:   cfg.Destination.IconEmoji,
		Text:        fmt.Sprintf("[%s] %s", alert.GetSeverityName(), message),
		Attachments: []SimpleAttachment{attachment},
	}

	return msg
}

func getSeverityColor(severity string) string {
	switch strings.ToLower(severity) {
	case "0", "1":
		return "#D62828" // Red
	case "2":
		return "#F77F00" // Orange
	case "3":
		return "#FFB703" // Yellow
	case "4", "5":
		return "#219EBC" // Blue
	case "6":
		return "#023047" // Dark Blue
	default:
		return "#999999" // Gray
	}
}
EOF

cat /tmp/notification_fixed.go
Output

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

// SimpleMessage for basic Mattermost/Slack compatibility
type SimpleMessage struct {
	Text        string                 `json:"text"`
	Channel     string                 `json:"channel,omitempty"`
	Username    string                 `json:"username,omitempty"`
	IconEmoji   string                 `json:"icon_emoji,omitempty"`
	Attachments []SimpleAttachment      `json:"attachments,omitempty"`
	Props       map[string]interface{} `json:"props,omitempty"`
}

// SimpleAttachment for message formatting
type SimpleAttachment struct {
	Fallback string                 `json:"fallback"`
	Color    string                 `json:"color"`
	Title    string                 `json:"title,omitempty"`
	Text     string                 `json:"text,omitempty"`
	Fields   []SimpleField           `json:"fields,omitempty"`
	Ts       int64                  `json:"ts,omitempty"`
}

// SimpleField for attachment fields
type SimpleField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// MessageClient handles posting to Slack or Mattermost
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

// PostMessage posts a message to Slack or Mattermost
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

// BuildMessage creates a formatted message for Slack or Mattermost
func BuildMessage(alert *GraylogAlert, cfg *Config) *SimpleMessage {
	color := getSeverityColor(alert.GetSeverity())
	timestamp := alert.GetTimestamp()
	message := alert.GetDisplayMessage()

	// Truncate message if too long
	if len(message) > 500 {
		message = message[:500] + "..."
	}

	// Build fields
	fields := []SimpleField{
		{
			Title: "Severity",
			Value: alert.GetSeverityName(),
			Short: true,
		},
		{
			Title: "Time",
			Value: timestamp.Format(time.RFC3339),
			Short: true,
		},
	}

	if alert.Source != "" {
		fields = append(fields, SimpleField{
			Title: "Source",
			Value: alert.Source,
			Short: true,
		})
	}

	if alert.EventDefinitionID != "" {
		fields = append(fields, SimpleField{
			Title: "Event ID",
			Value: alert.EventDefinitionID,
			Short: true,
		})
	}

	// Determine channel
	channel := cfg.Destination.Channel
	if cfg.Destination.Destinations != nil {
		if dest, exists := cfg.Destination.Destinations[alert.GetSeverity()]; exists {
			channel = dest
		}
	}

	// Build attachment
	attachment := SimpleAttachment{
		Fallback: fmt.Sprintf("[%s] %s", alert.GetSeverityName(), message),
		Color:    color,
		Title:    message,
		Fields:   fields,
		Ts:       timestamp.Unix(),
	}

	// Build message
	msg := &SimpleMessage{
		Channel:     channel,
		Username:    cfg.Destination.Username,
		IconEmoji:   cfg.Destination.IconEmoji,
		Text:        fmt.Sprintf("[%s] %s", alert.GetSeverityName(), message),
		Attachments: []SimpleAttachment{attachment},
	}

	return msg
}

func getSeverityColor(severity string) string {
	switch strings.ToLower(severity) {
	case "0", "1":
		return "#D62828" // Red
	case "2":
		return "#F77F00" // Orange
	case "3":
		return "#FFB703" // Yellow
	case "4", "5":
		return "#219EBC" // Blue
	case "6":
		return "#023047" // Dark Blue
	default:
		return "#999999" // Gray
	}
}