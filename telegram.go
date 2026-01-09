package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TelegramClient struct {
	BotToken string
	ChatID   string
	client   *http.Client
}

func NewTelegramClient(cfg *Config) *TelegramClient {
	return &TelegramClient{
		BotToken: cfg.BotToken,
		ChatID:   cfg.ChatID,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *TelegramClient) SendMessage(text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)

	payload := map[string]string{
		"chat_id": t.ChatID,
		"text":    text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	resp, err := t.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned %d", resp.StatusCode)
	}

	return nil
}
