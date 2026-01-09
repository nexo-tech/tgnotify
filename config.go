package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BotToken string `toml:"bot_token"`
	ChatID   string `toml:"chat_id"`
}

func LoadConfig() (*Config, error) {
	path := os.Getenv("TGNOTIFY_CONFIG")
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("get home dir: %w", err)
		}
		path = filepath.Join(home, ".tg-notify", "config.toml")
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("decode config %s: %w", path, err)
	}

	if cfg.BotToken == "" {
		return nil, fmt.Errorf("bot_token is required in %s", path)
	}
	if cfg.ChatID == "" {
		return nil, fmt.Errorf("chat_id is required in %s", path)
	}

	return &cfg, nil
}
