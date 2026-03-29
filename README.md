# zai-keycheck

Z.AI quota checker for GLM Coding Plan.

## About

Monitors your Z.AI API quota, checking TIME_LIMIT and TOKENS_LIMIT usage with automatic reset tracking. Supports multiple providers and stores API key encoded in base64 for security.

## Features

- ✅ Check TIME_LIMIT (requests per time window)
- ✅ Check TOKENS_LIMIT (token usage)
- ✅ Convert timestamps to your local timezone
- ✅ Multiple providers in the same config file
- ✅ Base64 API key for security (optional)
- ✅ Automatically updates `last_attempt` in config
- ✅ Default config location: `~/.config/zai-keycheck/`

## Installation

### Via Go (recommended)

```bash
go install github.com/victorhdchagas/zai-keycheck@latest
```

### Manual

```bash
git clone https://github.com/victorhdchagas/zai-keycheck.git
cd zai-keycheck
go build
mv zai-keycheck ~/.local/bin/  # or ~/go/bin/
```

## Usage

### First time setup

1. Encode your API key in base64:
```bash
zai-keycheck --encode YOUR_API_KEY
```

2. Create config file at `~/.config/zai-keycheck/providers.json`:

```json
{
  "api_key_base64": "your-base64-encoded-key",
  "providers": [
    {
      "url": "https://api.z.ai/api/monitor/usage/quota/limit",
      "available_at": "",
      "last_attempt": ""
    }
  ]
}
```

Or use plain `api_key`:
```json
{
  "api_key": "sk-your-key-here",
  "providers": [...]
}
```

### Check quota

```bash
zai-keycheck                    # uses default config
zai-keycheck -c ./config.json  # uses custom file
zai-keycheck --help             # shows help
```

### Example output

```
🔗 https://api.z.ai/api/monitor/usage/quota/limit

   📊 Plano: lite

   ⏱️  TIME_LIMIT (1 req / 5 min)
      Used: 100 | Remaining: 100 | 0%
      By model:
        - search-prime: 0
        - web-reader: 0
        - zread: 0
      🔄 Resets in: 112h 2m 30s (02/04/2026 22:37:15)

   🪙 TOKENS_LIMIT (5 req / 3 h)
      Usage: 1%
      🔄 Resets in: 3h 58m 14s (29/03/2026 10:32:59)

📁 Config updated: /home/user/.config/zai-keycheck/providers.json
```

## Security

The API key can be stored in two ways:

1. **api_key** — plain text (recommended for development only)
2. **api_key_base64** — base64 encoded (recommended for production/sharing)

Use `api_key_base64` when:
- Sharing config file
- Committing to Git
- Storing in backups

## License

MIT - see [LICENSE](LICENSE) file for details.
