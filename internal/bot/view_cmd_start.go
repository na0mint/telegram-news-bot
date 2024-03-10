package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/internal/botkit"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error {
		if _, err := api.Send(tgbotapi.NewMessage(update.FromChat().ID, "Hello world")); err != nil {
			return err
		}

		return nil
	}
}
