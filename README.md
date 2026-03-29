# zai_quotecheck

Z.AI quota checker for GLM Coding Plan.

## About

Monitors your Z.AI API quota, checking TIME_LIMIT and TOKENS_LIMIT usage with automatic reset tracking. Supports multiple providers and stores API key encoded in base64 for security.

## Features

- ✅ Check TIME_LIMIT (requests per time window)
- ✅ Check TOKENS_LIMIT (token usage)
- ✅ Convert timestamps to your local timezone
- ✅ Multiple providers in same config file
- ✅ Base64 API key for security (optional)
- ✅ Automatically updates `last_attempt` in config
- ✅ Default config location: `~/.config/zai_quotecheck/`
- ✅ `--init` command for quick config creation
- ✅ Docker support

## Installation

### Via Go (recommended)

```bash
go install github.com/victorhdchagas/zai_quotecheck@latest
```

### Manual

```bash
git clone https://github.com/victorhdchagas/zai_quotecheck.git
cd zai_quotecheck
go build
mv zai_quotecheck ~/.local/bin/  # or ~/go/bin/
```

### Via Docker

```bash
docker run --rm \
  -v ~/.config/zai_quotecheck:/root/.config/zai_quotecheck \
  -e TZ=America/Sao_Paulo \
  victorhdchagas/zai_quotecheck:latest
```

Or using docker-compose:
```bash
docker-compose run --rm zai_quotecheck
```

## Usage

### First time setup

**Option 1: Use --encode (recommended)**

```bash
zai_quotecheck --encode YOUR_API_KEY
```

This encodes your API key in base64 and saves it directly to the config file at `~/.config/zai_quotecheck/providers.json`.

**Option 2: Use --init**

```bash
zai_quotecheck --init
```

This creates an example config file at `~/.config/zai_quotecheck/providers.json`.

For manual editing, the config supports `api_key` (plain text) or `api_key_base64` (base64 encoded).

### Environment Variable (Priority over config file)

The API key can be set via environment variable, which has priority over config file:

```bash
export ZAI_API_KEY=your-api-key-here
zai_quotecheck
```

This is useful for Docker or CI/CD environments.

### Check quota

```bash
zai_quotecheck                    # uses default config
zai_quotecheck -c ./config.json  # uses custom file
zai_quotecheck --help             # shows help
```

### Check quota

```bash
zai_quotecheck                    # uses default config
zai_quotecheck -c ./config.json  # uses custom file
zai_quotecheck --help             # shows help
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

📁 Config updated: /home/user/.config/zai_quotecheck/providers.json
```

## Docker

### Build

```bash
docker build -t zai_quotecheck .
```

### Run

```bash
docker run --rm \
  -v ~/.config/zai_quotecheck:/root/.config/zai_quotecheck \
  -e ZAI_API_KEY=your-api-key-here \
  -e TZ=America/Sao_Paulo \
  zai_quotecheck
```

Or set via environment variable:
```bash
export ZAI_API_KEY=your-api-key-here
docker run --rm \
  -v ~/.config/zai_quotecheck:/root/.config/zai_quotecheck \
  -e ZAI_API_KEY \
  -e TZ=America/Sao_Paulo \
  zai_quotecheck
```

### Docker Compose

```bash
ZAI_API_KEY=your-api-key-here docker-compose run --rm zai_quotecheck
```

Or set in `.env` file (recommended):
```bash
echo "ZAI_API_KEY=your-api-key-here" > .env
docker-compose run --rm zai_quotecheck
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
