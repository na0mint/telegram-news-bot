package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"tg-bot/internal/botkit"
)

type SourceDeleter interface {
	Delete(ctx context.Context, id int64) error
}

func ViewCmdDeleteSource(deleter SourceDeleter) botkit.ViewFunc {
	return func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error {
		targetId, err := strconv.ParseInt(update.Message.CommandArguments(),
			10, 64)
		if err != nil {
			_, sendErr := api.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				"Failed to parse source id"))
			if sendErr != nil {
				return sendErr
			}

			return err
		}

		err = deleter.Delete(ctx, targetId)
		if err != nil {
			_, sendErr := api.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				"Failed to delete source"))
			if sendErr != nil {
				return sendErr
			}

			return err
		}

		_, sendErr := api.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Source with ID: `%d` succussfully deleted",
				targetId)))
		if sendErr != nil {
			return sendErr
		}

		return nil
	}
}
