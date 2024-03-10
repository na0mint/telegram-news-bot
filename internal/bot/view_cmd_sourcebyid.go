package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"tg-bot/internal/botkit"
	"tg-bot/internal/botkit/markup"
	"tg-bot/internal/model"
)

type SourceFinder interface {
	SourceById(ctx context.Context, id int64) (*model.Source, error)
}

func ViewCmdGetSourceById(finder SourceFinder) botkit.ViewFunc {
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

		source, err := finder.SourceById(ctx, targetId)

		var (
			msgText = fmt.Sprintf(
				"🌐 *%s*\nID: `%d`\nFeed URL: %s",
				markup.EscapeForMarkdown(source.Name),
				source.ID,
				markup.EscapeForMarkdown(source.FeedURL),
			)

			reply = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		reply.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := api.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
