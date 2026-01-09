package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HookPayload struct {
	SessionID        string `json:"session_id"`
	TranscriptPath   string `json:"transcript_path"`
	Cwd              string `json:"cwd"`
	HookEventName    string `json:"hook_event_name"`
	StopHookActive   bool   `json:"stop_hook_active"`
	Message          string `json:"message"`
	NotificationType string `json:"notification_type"`
}

type TranscriptEntry struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Message   struct {
		Role    string `json:"role"`
		Content any    `json:"content"`
	} `json:"message"`
}

type TranscriptContext struct {
	LastUserMessage string
	SessionDuration time.Duration
}

func ParseHookPayload(data []byte) (*HookPayload, error) {
	var payload HookPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func (h *HookPayload) ProjectName() string {
	return filepath.Base(h.Cwd)
}

func (h *HookPayload) GetTranscriptContext() *TranscriptContext {
	if h.TranscriptPath == "" {
		return nil
	}

	f, err := os.Open(h.TranscriptPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var ctx TranscriptContext
	var firstTimestamp, lastTimestamp time.Time
	scanner := bufio.NewScanner(f)

	// Increase buffer size for long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		var entry TranscriptEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		if firstTimestamp.IsZero() && !entry.Timestamp.IsZero() {
			firstTimestamp = entry.Timestamp
		}
		if !entry.Timestamp.IsZero() {
			lastTimestamp = entry.Timestamp
		}

		if entry.Type == "user" {
			ctx.LastUserMessage = extractContent(entry.Message.Content)
		}
	}

	if !firstTimestamp.IsZero() && !lastTimestamp.IsZero() {
		ctx.SessionDuration = lastTimestamp.Sub(firstTimestamp)
	}

	return &ctx
}

func extractContent(content any) string {
	switch v := content.(type) {
	case string:
		return truncate(v, 100)
	case []any:
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				if t, ok := m["type"].(string); ok && t == "text" {
					if text, ok := m["text"].(string); ok {
						return truncate(text, 100)
					}
				}
			}
		}
	}
	return ""
}

func truncate(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}
