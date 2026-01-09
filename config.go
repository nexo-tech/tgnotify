package main

import (
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
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".tgnotify.toml")
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
