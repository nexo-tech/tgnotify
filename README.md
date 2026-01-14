# tgnotify

Telegram notifications for [Claude Code](https://claude.ai/code) and [OpenCode](https://opencode.ai) hooks.

Get instant Telegram messages when your AI coding assistant finishes tasks or needs your attention.

## Features

- Notifications on task completion
- Alerts when assistant needs permission
- Session context (last task, duration)
- Simple TOML configuration
- Nix flake for easy installation
- Supports both Claude Code and OpenCode

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

### 4. Add Hooks

#### Claude Code

Add to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "Stop": [{ "hooks": [{ "type": "command", "command": "tgnotify" }] }],
    "Notification": [{ "hooks": [{ "type": "command", "command": "tgnotify" }] }]
  }
}
```

#### OpenCode

Create `~/.config/opencode/plugin/tgnotify.ts`:

```typescript
import { execSync } from "child_process";

export const TgNotifyPlugin = async ({ project }) => ({
  event: async ({ event }) => {
    if (event.type === "session.idle" || event.type === "session.error") {
      const payload = JSON.stringify({
        hook_event_name: event.type === "session.idle" ? "Stop" : "Notification",
        cwd: project.path,
        message: event.type === "session.error" ? "Session error" : "",
        tool_name: "opencode"
      });
      try {
        execSync("tgnotify", { input: payload, stdio: ["pipe", "inherit", "inherit"] });
      } catch {}
    }
  }
});
```

## Notification Examples

**Task Completed:**
```
｡✧ done! ✧｡

my-project
implement user authentication

/Users/me/projects/my-project
me@macbook · 5m 32s · 14:32
```

**Awaiting Input:**
```
｡･ﾟ need you~ ﾟ･｡

my-project
permission needed for: git push

/Users/me/projects/my-project
me@macbook · 14:35
```

## License

MIT
