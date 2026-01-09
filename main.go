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
	sys := GetSystemInfo()
	ctx := p.GetTranscriptContext()
	now := time.Now()
	cwd, _ := os.Getwd()

	// Header
	switch p.HookEventName {
	case "Stop":
		b.WriteString("｡✧ <b>done!</b> ✧｡\n\n")
	case "Notification":
		b.WriteString("｡･ﾟ <b>need you~</b> ﾟ･｡\n\n")
	default:
		b.WriteString(fmt.Sprintf("✿ <b>%s</b> ✿\n\n", strings.ToLower(p.HookEventName)))
	}

	// Project
	b.WriteString(fmt.Sprintf("<b>%s</b>\n", p.ProjectName()))

	// Task description
	if p.HookEventName == "Notification" && p.Message != "" {
		b.WriteString(fmt.Sprintf("<i>%s</i>\n", p.Message))
	}

	if p.HookEventName == "Stop" && ctx != nil && ctx.LastUserMessage != "" {
		b.WriteString(fmt.Sprintf("<i>%s</i>\n", ctx.LastUserMessage))
	}

	b.WriteString("\n")

	// Metadata - minimal
	b.WriteString(fmt.Sprintf("<code>%s</code>\n", cwd))
	b.WriteString(fmt.Sprintf("<code>%s@%s</code>", sys.Username, sys.Hostname))

	if ctx != nil && ctx.SessionDuration > 0 {
		b.WriteString(fmt.Sprintf(" · %s", formatDuration(ctx.SessionDuration)))
	}

	b.WriteString(fmt.Sprintf(" · %s", now.Format("15:04")))

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
