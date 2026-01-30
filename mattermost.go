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

type MattermostMessage struct {
	Channel     string       `json:"channel,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Fallback   string  `json:"fallback"`
	Color      string  `json:"color"`
	AuthorName string  `json:"author_name,omitempty"`
	Title      string  `json:"title,omitempty"`
	Text       string  `json:"text,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
	Timestamp  int64   `json:"ts,omitempty"`
}

type Field struct {
	Short bool   `json:"short"`
	Title string `json:"title"`
	Value string `json:"value"`
}

type MattermostClient struct {
	webhookURL string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewMattermostClient(webhookURL string, logger *zap.Logger) *MattermostClient {
	return &MattermostClient{
		webhookURL: webhookURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (mc *MattermostClient) PostMessage(msg *MattermostMessage) error {
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
		return fmt.Errorf("mattermost returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// BuildMessage creates a formatted message from a Graylog alert
func BuildMessage(alert *GraylogAlert, config *Config) *MattermostMessage {
	color := getSeverityColor(alert.GetSeverity())
	timestamp := alert.GetTimestamp()
	message := alert.GetDisplayMessage()

	// Truncate message if too long
	if len(message) > 500 {
		message = message[:500] + "..."
	}

	attachment := Attachment{
		Fallback:   fmt.Sprintf("[%s] %s", alert.GetSeverityName(), message),
		Color:      color,
		AuthorName: "Graylog",
		Title:      message,
		Timestamp:  timestamp.Unix(),
	}

	// Add fields
	attachment.Fields = append(attachment.Fields, Field{
		Short: true,
		Title: "Severity",
		Value: alert.GetSeverityName(),
	})

	attachment.Fields = append(attachment.Fields, Field{
		Short: true,
		Title: "Time",
		Value: timestamp.Format(time.RFC3339),
	})

	if alert.Source != "" {
		attachment.Fields = append(attachment.Fields, Field{
			Short: true,
			Title: "Source",
			Value: alert.Source,
		})
	}

	if alert.EventDefinitionID != "" {
		attachment.Fields = append(attachment.Fields, Field{
			Short: true,
			Title: "Event ID",
			Value: alert.EventDefinitionID,
		})
	}

	// Determine channel
	channel := config.Mattermost.Channel
	if config.Mattermost.Destinations != nil {
		if dest, exists := config.Mattermost.Destinations[alert.GetSeverity()]; exists {
			channel = dest
		}
	}

	msg := &MattermostMessage{
		Channel:     channel,
		Username:    config.Mattermost.Username,
		IconEmoji:   config.Mattermost.IconEmoji,
		Text:        fmt.Sprintf("[%s] %s", alert.GetSeverityName(), message),
		Attachments: []Attachment{attachment},
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
