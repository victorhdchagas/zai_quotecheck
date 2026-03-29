# zai_quotecheck

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
- ✅ Configuração padrão em `~/.config/zai_quotecheck/`
- ✅ Comando `--init` para criação rápida de config
- ✅ Suporte a Docker

## Instalação

### Via Go (recomendado)

```bash
go install github.com/victorhdchagas/zai_quotecheck@latest
```

### Manual

```bash
git clone https://github.com/victorhdchagas/zai_quotecheck.git
cd zai_quotecheck
go build
mv zai_quotecheck ~/.local/bin/  # ou ~/go/bin/
```

### Via Docker

```bash
docker run --rm \
  -v ~/.config/zai_quotecheck:/root/.config/zai_quotecheck \
  -e TZ=America/Sao_Paulo \
  victorhdchagas/zai_quotecheck:latest
```

Ou usando docker-compose:
```bash
docker-compose run --rm zai_quotecheck
```

## Uso

### Primeiro uso

**Opção 1: Usar --encode (recomendado)**

```bash
zai_quotecheck --encode SUA_API_KEY
```

Isso codifica sua API key em base64 e salva direto no arquivo de configuração em `~/.config/zai_quotecheck/providers.json`.

**Opção 2: Usar --init**

```bash
zai_quotecheck --init
```

Isso cria um arquivo de configuração de exemplo em `~/.config/zai_quotecheck/providers.json`.

Para editar manualmente, o arquivo suporta `api_key` (texto) ou `api_key_base64` (base64).

### Variável de ambiente (Prioridade sobre arquivo de config)

A API key pode ser definida via variável de ambiente, que tem prioridade sobre o arquivo de config:

```bash
export ZAI_API_KEY=sua-api-key-aqui
zai_quotecheck
```

Isso é útil para ambientes Docker ou CI/CD.

### Verificar quota

```bash
zai_quotecheck                    # usa config padrão
zai_quotecheck -c ./config.json  # usa arquivo customizado
zai_quotecheck --help             # mostra ajuda
```

### Verificar quota

```bash
zai_quotecheck                    # usa config padrão
zai_quotecheck -c ./config.json  # usa arquivo customizado
zai_quotecheck --help             # mostra ajuda
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

📁 Config atualizado: /home/wutachi/.config/zai_quotecheck/providers.json
```

## Docker

### Build

```bash
docker build -t zai_quotecheck .
```

### Executar

```bash
docker run --rm \
  -v ~/.config/zai_quotecheck:/root/.config/zai_quotecheck \
  -e ZAI_API_KEY=sua-api-key-aqui \
  -e TZ=America/Sao_Paulo \
  zai_quotecheck
```

Ou definir via variável de ambiente:
```bash
export ZAI_API_KEY=sua-api-key-aqui
docker run --rm \
  -v ~/.config/zai_quotecheck:/root/.config/zai_quotecheck \
  -e ZAI_API_KEY \
  -e TZ=America/Sao_Paulo \
  zai_quotecheck
```

### Docker Compose

```bash
ZAI_API_KEY=sua-api-key-aqui docker-compose run --rm zai_quotecheck
```

Ou criar arquivo `.env` (recomendado):
```bash
echo "ZAI_API_KEY=sua-api-key-aqui" > .env
docker-compose run --rm zai_quotecheck
```

### Executar

```bash
docker run --rm \
  -v ~/.config/zai_quotecheck:/root/.config/zai_quotecheck \
  -e TZ=America/Sao_Paulo \
  zai_quotecheck
```

### Docker Compose

```bash
docker-compose run --rm zai_quotecheck
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
