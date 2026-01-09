# 🤖 tgnotify

Telegram notifications for [Claude Code](https://claude.ai/code) hooks.

Get instant Telegram messages when Claude finishes tasks or needs your attention.

## Features

- 📬 Notifications on task completion
- ⚠️ Alerts when Claude needs permission
- 📊 Session context (last task, duration)
- 🔧 Simple TOML configuration
- ❄️ Nix flake for easy installation

## Installation

### Nix Flake (recommended)

```nix
{
  inputs.tgnotify.url = "github:nexo-tech/tgnotify";

  # Add to home.packages
  home.packages = [ inputs.tgnotify.packages.${system}.default ];
}
```

### Go Install

```bash
go install github.com/nexo-tech/tgnotify@latest
```

### Build from source

```bash
git clone https://github.com/nexo-tech/tgnotify
cd tgnotify
go build -o tgnotify
```

## Setup

### 1. Create Telegram Bot

1. Message [@BotFather](https://t.me/BotFather) on Telegram
2. Send `/newbot` and follow prompts
3. Copy the bot token

### 2. Get Chat ID

1. Message your new bot
2. Visit `https://api.telegram.org/bot<TOKEN>/getUpdates`
3. Find your `chat.id` in the response

### 3. Configure tgnotify

Create `~/.tg-notify/config.toml`:

```toml
bot_token = "123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
chat_id = "987654321"
```

Or set custom config path via environment variable:

```bash
export TGNOTIFY_CONFIG=/path/to/config.toml
```

### 4. Add Claude Code Hooks

Add to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "tgnotify"
          }
        ]
      }
    ],
    "Notification": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "tgnotify"
          }
        ]
      }
    ]
  }
}
```

## Notification Examples

**Task Completed:**
```
✅ Claude Code Completed

📁 my-project
📂 /Users/me/code/my-project
💬 "implement user authentication"
⏱️ Session: 5m 32s
```

**Needs Attention:**
```
⚠️ Claude Code Needs Attention

📁 my-project
📂 /Users/me/code/my-project
🔔 Claude needs your permission to use Bash
```

## License

MIT
