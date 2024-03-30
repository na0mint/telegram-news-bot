package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"log"
	"sync"
	"time"
)

type Config struct {
	TgBotToken              string        `hcl:"tg_bot_token" env:"TG_BOT_TOKEN" required:"true"`
	TgChannelId             int64         `hcl:"tg_channel_id" env:"TG_CHANNEL_ID" required:"true"`
	DbDSN                   string        `hcl:"db_dsn" env:"DB_DSN" default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	FetchInterval           time.Duration `hcl:"fetch_interval" env:"FETCH_INTERVAL" default:"100000m"`
	NotificationInterval    time.Duration `hcl:"notification_interval" env:"NOTIFICATION_INTERVAL" default:"1m"`
	FilterKeywords          []string      `hcl:"filter_keywords" env:"FILTER_KEYWORDS"`
	OpenAIKey               string        `hcl:"open_ai_key" env:"OPENAI_KEY"`
	OpenAIDefaultPrompt     string        `hcl:"open_ai_default_prompt" env:"OPENAI_DEFAULT_PROMPT"`
	OpenAITranslationPrompt string        `hcl:"open_ai_translation_prompt" env:"OPENAI_TRANSLATION_PROMPT"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "TGB",
			Files:     []string{"./config.hcl", "./config.local.hcl"},
			FileDecoders: map[string]aconfig.FileDecoder{
				".hcl": aconfighcl.New(),
			},
		})

		if err := loader.Load(); err != nil {
			log.Printf("[ERROR] failed to load config: %v", err)
		}
	})

	return cfg
}
