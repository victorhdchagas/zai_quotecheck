// zai_quotecheck - Z.AI quota checker for GLM Coding Plan
// Version: 1.0.0
// Author: Victor Chagas <wutachi@gmail.com>

package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const version = "1.0.0"

type Provider struct {
	URL         string `json:"url"`
	AvailableAt string `json:"available_at"`
	LastAttempt string `json:"last_attempt"`
}

type Config struct {
	APIKey       string     `json:"api_key,omitempty"`
	APIKeyBase64 string     `json:"api_key_base64,omitempty"`
	Providers    []Provider `json:"providers"`
}

type QuotaResponse struct {
	Code    int       `json:"code"`
	Msg     string    `json:"msg"`
	Success bool      `json:"success"`
	Data    QuotaData `json:"data"`
}

type QuotaData struct {
	Limits []Limit `json:"limits"`
	Level  string  `json:"level"`
}

type Limit struct {
	Type          string        `json:"type"`
	Unit          int           `json:"unit"`
	Number        int           `json:"number"`
	Usage         int           `json:"usage"`
	CurrentValue  int           `json:"currentValue"`
	Remaining     int           `json:"remaining"`
	Percentage    float64       `json:"percentage"`
	NextResetTime int64         `json:"nextResetTime"`
	UsageDetails  []UsageDetail `json:"usageDetails"`
}

type UsageDetail struct {
	ModelCode string `json:"modelCode"`
	Usage     int    `json:"usage"`
}

func getConfigPath() string {
	configPath := *flag.String("config", "", "Caminho do arquivo de configuração")
	if configPath == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao obter diretório de config: %v\n", err)
			os.Exit(1)
		}
		appDir := filepath.Join(configDir, "zai_quotecheck")
		if err := os.MkdirAll(appDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar diretório %s: %v\n", appDir, err)
			os.Exit(1)
		}
		return filepath.Join(appDir, "providers.json")
	}
	return configPath
}

func loadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(filePath string, config *Config) error {
	out, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, out, 0644)
}

func encodeToConfig(key string) {
	filePath := getConfigPath()

	config, err := loadConfig(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler config: %v\n", err)
		os.Exit(1)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(key))
	config.APIKeyBase64 = encoded
	config.APIKey = ""

	if err := saveConfig(filePath, config); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao salvar config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ API key codificada e salva em: %s\n", filePath)
	fmt.Printf("   Base64: %s\n", encoded)
}

func main() {
	help := flag.Bool("help", false, "Mostra esta ajuda")
	encodeKey := flag.String("encode", "", "Codifica uma API key e salva na config")
	initConfig := flag.Bool("init", false, "Cria arquivo de configuração de exemplo")

	flag.Parse()

	if *encodeKey != "" {
		encodeToConfig(*encodeKey)
		return
	}

	if *help {
		printHelp()
		return
	}

	if *initConfig {
		createExampleConfig()
		return
	}

	var filePath string
	var config *Config

	apiKey := os.Getenv("ZAI_API_KEY")
	if apiKey == "" {
		filePath = getConfigPath()

		var err error
		config, err = loadConfig(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao ler %s: %v\n", filePath, err)
			fmt.Fprintf(os.Stderr, "\nUse --init para criar arquivo de exemplo, ou --config para especificar outro caminho.\n")
			os.Exit(1)
		}

		apiKey = config.APIKey
		if apiKey == "" && config.APIKeyBase64 != "" {
			decoded, err := base64.StdEncoding.DecodeString(config.APIKeyBase64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao decodificar api_key_base64: %v\n", err)
				os.Exit(1)
			}
			apiKey = string(decoded)
		}

		if apiKey == "" {
			fmt.Fprintln(os.Stderr, "api_key ou api_key_base64 não definida no JSON.")
			fmt.Fprintln(os.Stderr, "Use: zai_quotecheck --encode SUA_API_KEY")
			fmt.Fprintln(os.Stderr, "Ou defina ZAI_API_KEY como variável de ambiente.")
			os.Exit(1)
		}
	}

	if apiKey != "" {
		config = &Config{
			Providers: []Provider{
				{
					URL:         "https://api.z.ai/api/monitor/usage/quota/limit",
					AvailableAt: "",
					LastAttempt: "",
				},
			},
		}
	}

	if len(config.Providers) == 0 {
		fmt.Println("Nenhum provider encontrado no arquivo.")
		return
	}

	local := time.Now().Location()
	now := time.Now()
	updated := false

	for i, p := range config.Providers {
		if i > 0 {
			fmt.Println(strings.Repeat("─", 60))
		}

		fmt.Printf("🔗 %s\n\n", p.URL)

		resp, err := checkQuota(p.URL, apiKey)
		if err != nil {
			fmt.Printf("   ❌ Erro na requisição: %v\n", err)
			continue
		}

		config.Providers[i].LastAttempt = now.UTC().Format(time.RFC3339)
		updated = true

		if !resp.Success {
			fmt.Printf("   ❌ Falha: %s (code: %d)\n", resp.Msg, resp.Code)
			continue
		}

		fmt.Printf("   📊 Plano: %s\n\n", resp.Data.Level)

		for _, limit := range resp.Data.Limits {
			printLimit(limit, local, now)
			fmt.Println()
		}
	}

	if updated && filePath != "" {
		if err := saveConfig(filePath, config); err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao salvar config: %v\n", err)
		} else {
			fmt.Printf("📁 Config atualizado: %s\n", filePath)
		}
	}
}

func createExampleConfig() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao obter diretório de config: %v\n", err)
		os.Exit(1)
	}
	appDir := filepath.Join(configDir, "zai_quotecheck")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar diretório %s: %v\n", appDir, err)
		os.Exit(1)
	}
	filePath := filepath.Join(appDir, "providers.json")

	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("Arquivo de configuração já existe em: %s\n", filePath)
		fmt.Println("Para recriar, apague o arquivo primeiro ou use --config para especificar outro caminho.")
		return
	}

	exampleConfig := Config{
		APIKeyBase64: "",
		Providers: []Provider{
			{
				URL:         "https://api.z.ai/api/monitor/usage/quota/limit",
				AvailableAt: "",
				LastAttempt: "",
			},
		},
	}

	if err := saveConfig(filePath, &exampleConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao salvar config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Arquivo de configuração criado: %s\n", filePath)
	fmt.Println("\nPróximos passos:")
	fmt.Println("  1. Encode sua API key: zai_quotecheck --encode SUA_API_KEY")
	fmt.Println("  2. Execute: zai_quotecheck")
}

func printHelp() {
	fmt.Printf("zai_quotecheck v%s - Verificador de quota Z.AI para GLM Coding Plan\n", version)
	fmt.Println("\nUSO:")
	fmt.Println("  zai_quotecheck [opções]\n")
	fmt.Println("OPÇÕES:")
	fmt.Println("  -config string")
	fmt.Println("      Caminho do arquivo de configuração (padrão: ~/.config/zai_quotecheck/providers.json)")
	fmt.Println("  -encode string")
	fmt.Println("      Codifica uma API key e salva na config")
	fmt.Println("  -init")
	fmt.Println("      Cria arquivo de configuração de exemplo")
	fmt.Println("  -help")
	fmt.Println("      Mostra esta ajuda\n")
	fmt.Println("SUBCOMANDOS:")
	fmt.Println("  zai_quotecheck --init")
	fmt.Println("      Cria arquivo de configuração padrão em ~/.config/zai_quotecheck/providers.json")
	fmt.Println("  zai_quotecheck --encode SUA_API_KEY")
	fmt.Println("      Codifica a API key em base64 e salva no arquivo de config\n")
	fmt.Println("ARQUIVO DE CONFIGURAÇÃO:")
	fmt.Println("  Padrão: ~/.config/zai_quotecheck/providers.json")
	fmt.Println("  API key pode ser definida de 3 formas (em ordem de prioridade):")
	fmt.Println("    1. Variável de ambiente: ZAI_API_KEY")
	fmt.Println("    2. api_key (texto no arquivo de config)")
	fmt.Println("    3. api_key_base64 (base64 no arquivo de config)")
	fmt.Println(`{
  "api_key_base64": "MTQ5NjNiYjMzMWI4NGJkNzg0ZTcxYWM3NzMxY2MzNzEuVm05bEF2Y1ZrZldVTGhoYw==",
  "providers": [
    {
      "url": "https://api.z.ai/api/monitor/usage/quota/limit",
      "available_at": "",
      "last_attempt": ""
    }
  ]
}`)
	fmt.Println("\nEXEMPLOS:")
	fmt.Println("  zai_quotecheck --init")
	fmt.Println("  zai_quotecheck --encode sk-sua-chave-aqui")
	fmt.Println("  zai_quotecheck")
	fmt.Println("  zai_quotecheck --config ./meu-config.json")
}

func printLimit(limit Limit, loc *time.Location, now time.Time) {
	switch limit.Type {
	case "TIME_LIMIT":
		fmt.Printf("   ⏱️  TIME_LIMIT (%d req / %d min)\n", limit.Number, limit.Unit)
		fmt.Printf("      Usadas: %d | Restantes: %d | %.0f%%\n",
			limit.Usage, limit.Remaining, limit.Percentage)
		if limit.CurrentValue > 0 {
			fmt.Printf("      Valor atual: %d\n", limit.CurrentValue)
		}
		if len(limit.UsageDetails) > 0 {
			fmt.Println("      Por modelo:")
			for _, d := range limit.UsageDetails {
				fmt.Printf("        - %s: %d\n", d.ModelCode, d.Usage)
			}
		}

	case "TOKENS_LIMIT":
		fmt.Printf("   🪙 TOKENS_LIMIT (%d req / %d h)\n", limit.Number, limit.Unit)
		fmt.Printf("      Uso: %.0f%%\n", limit.Percentage)
	}

	if limit.NextResetTime > 0 {
		resetTime := time.Unix(limit.NextResetTime/1000, 0).In(loc)
		if resetTime.After(now) {
			d := resetTime.Sub(now)
			fmt.Printf("      🔄 Reseta em: %s (%s)\n",
				formatDuration(d),
				resetTime.Format("02/01/2006 15:04:05"))
		} else {
			fmt.Println("      ✅ Já resetou!")
		}
	}
}

func checkQuota(url, apiKey string) (*QuotaResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept-Language", "en-US,en")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var quota QuotaResponse
	if err := json.Unmarshal(body, &quota); err != nil {
		return nil, fmt.Errorf("erro ao decodificar: %v (body: %s)", err, string(body))
	}

	return &quota, nil
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%d segundos", int(d.Seconds()))
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	parts := []string{}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}
