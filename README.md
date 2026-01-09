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

Create `~/.tgnotify.toml`:

```toml
bot_token = "123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
chat_id = "987654321"
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
━━━━━━━━━━━━━━━━━━━━━━
✅  TASK COMPLETED
━━━━━━━━━━━━━━━━━━━━━━

📁  my-project

💬  "implement user authentication"

┌─────────────────────
│ 📂  /Users/me/projects/my-project
│ 👤  me@macbook
│ 🔑  abc12345
│ ⏱   5m 32s
│ 🕐  14:32:07
└─────────────────────
```

**Awaiting Input:**
```
━━━━━━━━━━━━━━━━━━━━━━
⏳  AWAITING INPUT
━━━━━━━━━━━━━━━━━━━━━━

📁  my-project

💬  Claude needs permission to execute: rm -rf node_modules

┌─────────────────────
│ 📂  /Users/me/projects/my-project
│ 👤  me@macbook
│ 🔑  xyz98765
│ 🕐  14:35:22
└─────────────────────
```

## License

MIT
