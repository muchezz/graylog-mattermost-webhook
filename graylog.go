package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GraylogAlert represents the webhook payload from Graylog
type GraylogAlert struct {
	EventDefinitionID   string                 `json:"event_definition_id"`
	EventDefinitionType string                 `json:"event_definition_type"`
	EventTriggerID      string                 `json:"event_trigger_id"`
	EventTimestamp      string                 `json:"event_timestamp"`
	EventMessage        string                 `json:"event_message"`
	Priority            string                 `json:"priority"`
	Alert               bool                   `json:"alert"`
	Fields              map[string]interface{} `json:"fields"`
	Source              string                 `json:"source"`
	Message             string                 `json:"message"`
	Timestamp           string                 `json:"timestamp"`
	Level               interface{}            `json:"level"`
	FullMessage         string                 `json:"full_message"`
}

// ParseGraylogAlert parses the JSON webhook payload
func ParseGraylogAlert(data []byte) (*GraylogAlert, error) {
	var alert GraylogAlert
	if err := json.Unmarshal(data, &alert); err != nil {
		return nil, fmt.Errorf("failed to parse Graylog alert: %w", err)
	}
	return &alert, nil
}

// GetSeverity returns the severity level
func (ga *GraylogAlert) GetSeverity() string {
	if ga.Priority != "" {
		return ga.Priority
	}
	if level, ok := ga.Level.(float64); ok {
		return fmt.Sprintf("%d", int(level))
	}
	return "unknown"
}

// GetSeverityName returns human-readable severity name
func (ga *GraylogAlert) GetSeverityName() string {
	switch strings.ToLower(ga.GetSeverity()) {
	case "0":
		return "EMERGENCY"
	case "1":
		return "ALERT"
	case "2":
		return "ERROR"
	case "3":
		return "WARNING"
	case "4":
		return "NOTICE"
	case "5":
		return "INFO"
	case "6":
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// GetTimestamp returns the alert timestamp
func (ga *GraylogAlert) GetTimestamp() time.Time {
	if ga.EventTimestamp != "" {
		if t, err := time.Parse(time.RFC3339, ga.EventTimestamp); err == nil {
			return t
		}
	}
	if ga.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, ga.Timestamp); err == nil {
			return t
		}
	}
	return time.Now()
}

// GetDisplayMessage returns the most relevant message
func (ga *GraylogAlert) GetDisplayMessage() string {
	if ga.EventMessage != "" {
		return ga.EventMessage
	}
	if ga.Message != "" {
		return ga.Message
	}
	if ga.FullMessage != "" {
		return ga.FullMessage
	}
	return "No message available"
}
