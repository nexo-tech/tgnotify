package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "tgnotify: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("read stdin: %w", err)
	}

	payload, err := ParseHookPayload(input)
	if err != nil {
		return fmt.Errorf("parse hook payload: %w", err)
	}

	msg := formatMessage(payload)

	client := NewTelegramClient(cfg)
	if err := client.SendMessage(msg); err != nil {
		return fmt.Errorf("send telegram: %w", err)
	}

	return nil
}

func formatMessage(p *HookPayload) string {
	var b strings.Builder

	switch p.HookEventName {
	case "Stop":
		b.WriteString("✅ Claude Code Completed\n\n")
	case "Notification":
		b.WriteString("⚠️ Claude Code Needs Attention\n\n")
	default:
		b.WriteString(fmt.Sprintf("🔔 Claude Code: %s\n\n", p.HookEventName))
	}

	b.WriteString(fmt.Sprintf("📁 %s\n", p.ProjectName()))
	b.WriteString(fmt.Sprintf("📂 %s\n", p.Cwd))

	if p.HookEventName == "Notification" && p.Message != "" {
		b.WriteString(fmt.Sprintf("🔔 %s\n", p.Message))
	}

	if p.HookEventName == "Stop" {
		if ctx := p.GetTranscriptContext(); ctx != nil {
			if ctx.LastUserMessage != "" {
				b.WriteString(fmt.Sprintf("💬 \"%s\"\n", ctx.LastUserMessage))
			}
			if ctx.SessionDuration > 0 {
				b.WriteString(fmt.Sprintf("⏱️ Session: %s\n", formatDuration(ctx.SessionDuration)))
			}
		}
	}

	return b.String()
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
