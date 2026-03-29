# zai-keycheck

Verificador de quota da Z.AI para GLM Coding Plan.

## Sobre

Monitora sua quota da API Z.AI, verificando o uso de TIME_LIMIT e TOKENS_LIMIT com reset automático. Suporta múltiplos providers e armazena a API key codificada em base64 para segurança.

## Recursos

- ✅ Verifica TIME_LIMIT (requisições por janela de tempo)
- ✅ Verifica TOKENS_LIMIT (limite de tokens)
- ✅ Converte timestamps para sua timezone local
- ✅ Múltiplos providers no mesmo arquivo de config
- ✅ API key em base64 para segurança (opcional)
- ✅ Atualiza automaticamente `last_attempt` no config
- ✅ Configuração padrão em `~/.config/zai-keycheck/`
- ✅ Comando `--init` para criação rápida de config
- ✅ Suporte a Docker

## Instalação

### Via Go (recomendado)

```bash
go install github.com/victorhdchagas/zai-keycheck@latest
```

### Manual

```bash
git clone https://github.com/victorhdchagas/zai-keycheck.git
cd zai-keycheck
go build
mv zai-keycheck ~/.local/bin/  # ou ~/go/bin/
```

### Via Docker

```bash
docker run --rm \
  -v ~/.config/zai-keycheck:/root/.config/zai-keycheck \
  -e TZ=America/Sao_Paulo \
  victorhdchagas/zai-keycheck:latest
```

Ou usando docker-compose:
```bash
docker-compose run --rm zai-keycheck
```

## Uso

### Primeiro uso

**Opção 1: Usar --encode (recomendado)**

```bash
zai-keycheck --encode SUA_API_KEY
```

Isso codifica sua API key em base64 e salva direto no arquivo de configuração em `~/.config/zai-keycheck/providers.json`.

**Opção 2: Usar --init**

```bash
zai-keycheck --init
```

Isso cria um arquivo de configuração de exemplo em `~/.config/zai-keycheck/providers.json`.

Para editar manualmente, o arquivo suporta `api_key` (texto) ou `api_key_base64` (base64).
```

### Verificar quota

```bash
zai-keycheck                    # usa config padrão
zai-keycheck -c ./config.json  # usa arquivo customizado
zai-keycheck --help             # mostra ajuda
```

### Saída de exemplo

```
🔗 https://api.z.ai/api/monitor/usage/quota/limit

   📊 Plano: lite

   ⏱️  TIME_LIMIT (1 req / 5 min)
      Usadas: 100 | Restantes: 100 | 0%
      Por modelo:
        - search-prime: 0
        - web-reader: 0
        - zread: 0
      🔄 Reseta em: 112h 2m 30s (02/04/2026 22:37:15)

   🪙 TOKENS_LIMIT (5 req / 3 h)
      Uso: 1%
      🔄 Reseta em: 3h 58m 14s (29/03/2026 10:32:59)

📁 Config atualizado: /home/wutachi/.config/zai-keycheck/providers.json
```

## Docker

### Build

```bash
docker build -t zai-keycheck .
```

### Executar

```bash
docker run --rm \
  -v ~/.config/zai-keycheck:/root/.config/zai-keycheck \
  -e TZ=America/Sao_Paulo \
  zai-keycheck
```

### Docker Compose

```bash
docker-compose run --rm zai-keycheck
```

## Segurança

A API key pode ser armazenada de duas formas:

1. **api_key** — texto plano (recomendado apenas para desenvolvimento)
2. **api_key_base64** — codificada em base64 (recomendado para produção/compartilhamento)

Use `api_key_base64` quando:
- Compartilhar o arquivo de config
- Commitar no Git
- Armazenar em backup

## Licença

MIT - ver arquivo [LICENSE](LICENSE) para detalhes.
