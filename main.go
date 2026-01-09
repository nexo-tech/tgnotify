package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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
	doneHeaders = []string{"｡✧ <b>done!</b> ✧｡", "✿ <b>finished~</b> ✿", "♪ <b>complete!</b> ♪", "☆ <b>all done!</b> ☆", "･ﾟ✧ <b>yay!</b> ✧ﾟ･", "♡ <b>done!</b> ♡"}
	needHeaders = []string{"｡･ﾟ <b>need you~</b> ﾟ･｡", "✿ <b>waiting...</b> ✿", "♪ <b>hey~</b> ♪", "☆ <b>psst!</b> ☆", "･ﾟ✧ <b>help?</b> ✧ﾟ･", "♡ <b>um...</b> ♡"}
)

func main() {
	path := os.Getenv("TGNOTIFY_CONFIG")
	if path == "" {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".tgnotify.toml")
	}

	var cfg struct {
		Token  string `toml:"bot_token"`
		ChatID string `toml:"chat_id"`
	}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return
	}

	var h struct {
		TranscriptPath string `json:"transcript_path"`
		Cwd            string `json:"cwd"`
		Event          string `json:"hook_event_name"`
		Message        string `json:"message"`
	}
	if json.NewDecoder(os.Stdin).Decode(&h) != nil {
		return
	}

	var b strings.Builder
	u, _ := user.Current()
	cwd, _ := os.Getwd()
	host, _ := os.Hostname()

	switch h.Event {
	case "Stop":
		b.WriteString(doneHeaders[rand.Intn(len(doneHeaders))])
	case "Notification":
		b.WriteString(needHeaders[rand.Intn(len(needHeaders))])
	default:
		fmt.Fprintf(&b, "✿ <b>%s</b> ✿", strings.ToLower(h.Event))
	}

	fmt.Fprintf(&b, "\n\n<b>%s</b>\n", filepath.Base(h.Cwd))

	if h.Event == "Notification" && h.Message != "" {
		fmt.Fprintf(&b, "<i>%s</i>\n", h.Message)
	}

	var dur time.Duration
	if h.Event == "Stop" {
		if msg, d := readTranscript(h.TranscriptPath); msg != "" {
			fmt.Fprintf(&b, "<i>%s</i>\n", msg)
			dur = d
		}
	}

	fmt.Fprintf(&b, "\n<code>%s</code>\n<code>%s@%s</code>", cwd, u.Username, host)
	if dur > 0 {
		fmt.Fprintf(&b, " · %s", fmtDur(dur))
	}
	fmt.Fprintf(&b, " · %s", time.Now().Format("15:04"))

	body, _ := json.Marshal(map[string]string{"chat_id": cfg.ChatID, "text": b.String(), "parse_mode": "HTML"})
	http.Post("https://api.telegram.org/bot"+cfg.Token+"/sendMessage", "application/json", bytes.NewReader(body))
}

func readTranscript(path string) (string, time.Duration) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0
	}
	defer f.Close()

	var msg string
	var first, last time.Time
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)

	for sc.Scan() {
		var e struct {
			Type    string    `json:"type"`
			Time    time.Time `json:"timestamp"`
			Message struct{ Content any } `json:"message"`
		}
		if json.Unmarshal(sc.Bytes(), &e) != nil {
			continue
		}
		if !e.Time.IsZero() {
			if first.IsZero() {
				first = e.Time
			}
			last = e.Time
		}
		if e.Type == "user" {
			if s, ok := e.Message.Content.(string); ok {
				msg = s
			} else if arr, ok := e.Message.Content.([]any); ok {
				for _, item := range arr {
					if m, ok := item.(map[string]any); ok && m["type"] == "text" {
						if t, ok := m["text"].(string); ok {
							msg = t
						}
					}
				}
			}
		}
	}

	msg = strings.TrimSpace(strings.ReplaceAll(msg, "\n", " "))
	if len(msg) > 100 {
		msg = msg[:100] + "..."
	}
	return msg, last.Sub(first)
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
