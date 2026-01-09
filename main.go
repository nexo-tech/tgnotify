package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	doneHeaders = []string{
		"｡✧ <b>done!</b> ✧｡", "✿ <b>finished~</b> ✿", "♪ <b>complete!</b> ♪",
		"☆ <b>all done!</b> ☆", "･ﾟ✧ <b>yay!</b> ✧ﾟ･", "♡ <b>done!</b> ♡",
	}
	needHeaders = []string{
		"｡･ﾟ <b>need you~</b> ﾟ･｡", "✿ <b>waiting...</b> ✿", "♪ <b>hey~</b> ♪",
		"☆ <b>psst!</b> ☆", "･ﾟ✧ <b>help?</b> ✧ﾟ･", "♡ <b>um...</b> ♡",
	}
)

type Config struct {
	BotToken string `toml:"bot_token"`
	ChatID   string `toml:"chat_id"`
}

type Hook struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	Cwd            string `json:"cwd"`
	Event          string `json:"hook_event_name"`
	Message        string `json:"message"`
}

func main() {
	path := os.Getenv("TGNOTIFY_CONFIG")
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".tgnotify.toml")
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return // silent exit if no config
	}

	input, _ := io.ReadAll(os.Stdin)
	var hook Hook
	if json.Unmarshal(input, &hook) != nil {
		return
	}

	msg := formatMsg(&hook)
	send(cfg.BotToken, cfg.ChatID, msg)
}

func formatMsg(h *Hook) string {
	var b strings.Builder
	cwd, _ := os.Getwd()
	user, _ := user.Current()
	host, _ := os.Hostname()

	// Header
	switch h.Event {
	case "Stop":
		b.WriteString(doneHeaders[rand.Intn(len(doneHeaders))] + "\n\n")
	case "Notification":
		b.WriteString(needHeaders[rand.Intn(len(needHeaders))] + "\n\n")
	default:
		fmt.Fprintf(&b, "✿ <b>%s</b> ✿\n\n", strings.ToLower(h.Event))
	}

	// Project & message
	fmt.Fprintf(&b, "<b>%s</b>\n", filepath.Base(h.Cwd))
	if h.Event == "Notification" && h.Message != "" {
		fmt.Fprintf(&b, "<i>%s</i>\n", h.Message)
	}
	if h.Event == "Stop" {
		if msg, dur := readTranscript(h.TranscriptPath); msg != "" {
			fmt.Fprintf(&b, "<i>%s</i>\n", msg)
			if dur > 0 {
				b.WriteString("\n<code>" + cwd + "</code>\n")
				fmt.Fprintf(&b, "<code>%s@%s</code> · %s · %s", user.Username, host, fmtDur(dur), time.Now().Format("15:04"))
				return b.String()
			}
		}
	}

	b.WriteString("\n<code>" + cwd + "</code>\n")
	fmt.Fprintf(&b, "<code>%s@%s</code> · %s", user.Username, host, time.Now().Format("15:04"))
	return b.String()
}

func readTranscript(path string) (lastMsg string, duration time.Duration) {
	if path == "" {
		return
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	var first, last time.Time
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		var e struct {
			Type      string    `json:"type"`
			Timestamp time.Time `json:"timestamp"`
			Message   struct {
				Content any `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal(scanner.Bytes(), &e) != nil {
			continue
		}
		if first.IsZero() && !e.Timestamp.IsZero() {
			first = e.Timestamp
		}
		if !e.Timestamp.IsZero() {
			last = e.Timestamp
		}
		if e.Type == "user" {
			lastMsg = extractText(e.Message.Content)
		}
	}
	if !first.IsZero() && !last.IsZero() {
		duration = last.Sub(first)
	}
	return
}

func extractText(c any) string {
	switch v := c.(type) {
	case string:
		return trunc(v)
	case []any:
		for _, item := range v {
			if m, ok := item.(map[string]any); ok && m["type"] == "text" {
				if t, ok := m["text"].(string); ok {
					return trunc(t)
				}
			}
		}
	}
	return ""
}

func trunc(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	if len(s) > 100 {
		return s[:100] + "..."
	}
	return s
}

func fmtDur(d time.Duration) string {
	d = d.Round(time.Second)
	if h := d / time.Hour; h > 0 {
		return fmt.Sprintf("%dh%dm", h, (d%time.Hour)/time.Minute)
	}
	if m := d / time.Minute; m > 0 {
		return fmt.Sprintf("%dm%ds", m, (d%time.Minute)/time.Second)
	}
	return fmt.Sprintf("%ds", d/time.Second)
}

func send(token, chatID, text string) {
	body, _ := json.Marshal(map[string]string{"chat_id": chatID, "text": text, "parse_mode": "HTML"})
	http.Post("https://api.telegram.org/bot"+token+"/sendMessage", "application/json", bytes.NewReader(body))
}
