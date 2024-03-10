package main

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tg-bot/internal/bot"
	"tg-bot/internal/bot/middleware"
	"tg-bot/internal/botkit"
	"tg-bot/internal/config"
	"tg-bot/internal/fetcher"
	"tg-bot/internal/notifier"
	"tg-bot/internal/storage"
	"tg-bot/internal/summary"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TgBotToken)
	if err != nil {
		log.Printf("[ERROR] Failed to create a bot: %v", err)
		return
	}

	db, err := sqlx.Connect("postgres", config.Get().DbDSN)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		postFetcher    = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		tgNotifier = notifier.NewNotifier(
			articleStorage,
			summary.NewOpenAISummarizer(config.Get().OpenAIKey, config.Get().OpenAIPrompt),
			botAPI,
			config.Get().NotificationInterval,
			2*config.Get().FetchInterval,
			config.Get().TgChannelId)
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	newsBot := botkit.NewBot(botAPI)
	newsBot.RegisterCmdView("start", bot.ViewCmdStart())
	newsBot.RegisterCmdView("addsource",
		middleware.AdminOnly(config.Get().TgChannelId,
			bot.ViewCmdAddSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView("sources",
		middleware.AdminOnly(config.Get().TgChannelId,
			bot.ViewCmdListSources(sourceStorage),
		),
	)
	newsBot.RegisterCmdView("sourcebyid",
		middleware.AdminOnly(config.Get().TgChannelId,
			bot.ViewCmdGetSourceById(sourceStorage),
		),
	)
	newsBot.RegisterCmdView("deletesource",
		middleware.AdminOnly(config.Get().TgChannelId,
			bot.ViewCmdDeleteSource(sourceStorage),
		),
	)

	go func(ctx context.Context) {
		if err := postFetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] Failed to start fetcher: %v", err)
				return
			}

			log.Println("fetcher stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := tgNotifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to start notifier: %v", err)
				return
			}

			log.Println("notifier stopped")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to start bot: %v", err)
			return
		}

		log.Println("bot stopped")
	}
}
