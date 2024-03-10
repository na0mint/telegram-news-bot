package middleware

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"tg-bot/internal/botkit"
)

func AdminOnly(channelID int64, next botkit.ViewFunc) botkit.ViewFunc {
	return func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error {
		admins, err := api.GetChatAdministrators(
			tgbotapi.ChatAdministratorsConfig{
				ChatConfig: tgbotapi.ChatConfig{
					ChatID: channelID,
				},
			},
		)

		if err != nil {
			return err
		}

		for _, admin := range admins {
			if admin.User.ID == update.Message.From.ID {
				return next(ctx, api, update)
			}
		}

		if _, err := api.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"You are not permitted to use this command",
		)); err != nil {
			log.Printf("[ERROR] Failed to send a message via telegram")
			return nil
		}

		return nil
	}
}
