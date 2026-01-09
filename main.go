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

	// Header
	switch p.HookEventName {
	case "Stop":
		b.WriteString("✅✅✅ <b>TASK COMPLETED</b> ✅✅✅\n\n")
	case "Notification":
		b.WriteString("⏳⏳⏳ <b>AWAITING INPUT</b> ⏳⏳⏳\n\n")
	default:
		b.WriteString(fmt.Sprintf("🔔🔔🔔 <b>%s</b> 🔔🔔🔔\n\n", strings.ToUpper(p.HookEventName)))
	}

	// Project
	b.WriteString(fmt.Sprintf("📦 <b>%s</b>\n\n", p.ProjectName()))

	// Task description for notifications
	if p.HookEventName == "Notification" && p.Message != "" {
		b.WriteString(fmt.Sprintf("💬 <i>%s</i>\n\n", p.Message))
	}

	// Last user message for completed tasks
	if p.HookEventName == "Stop" && ctx != nil && ctx.LastUserMessage != "" {
		b.WriteString(fmt.Sprintf("💬 <i>\"%s\"</i>\n\n", ctx.LastUserMessage))
	}

	// Metadata
	b.WriteString("📂 <code>" + p.Cwd + "</code>\n")

	if sys.Username != "" && sys.Hostname != "" {
		b.WriteString(fmt.Sprintf("👤 <code>%s@%s</code>\n", sys.Username, sys.Hostname))
	}

	if p.SessionID != "" {
		sid := p.SessionID
		if len(sid) > 8 {
			sid = sid[:8]
		}
		b.WriteString(fmt.Sprintf("🔑 <code>%s</code>\n", sid))
	}

	if ctx != nil && ctx.SessionDuration > 0 {
		b.WriteString(fmt.Sprintf("⏱️ <code>%s</code>\n", formatDuration(ctx.SessionDuration)))
	}

	b.WriteString(fmt.Sprintf("🕐 <code>%s</code>\n", now.Format("15:04:05")))

	b.WriteString("\n🤖🤖🤖🤖🤖🤖🤖🤖🤖🤖")

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
