package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/internal/botkit"
	"tg-bot/internal/model"
)

type SourceStorage interface {
	Save(ctx context.Context, source model.Source) (int64, error)
}

func ViewCmdAddSource(storage SourceStorage) botkit.ViewFunc {
	type addSourceArgs struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		TopicID int64  `json:"topicID"`
		Type    string `json:"type"`
	}

	return func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[addSourceArgs](update.Message.CommandArguments())
		if err != nil {
			_, sendErr := api.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				"Failed to parse command arguments"))
			if sendErr != nil {
				return sendErr
			}

			return err
		}

		source := model.Source{
			Name:    args.Name,
			FeedURL: args.URL,
			TopicID: args.TopicID,
			Type:    args.Type,
		}

		sourceID, err := storage.Save(ctx, source)
		if err != nil {
			return err
		}

		var (
			msgText = fmt.Sprintf(
				"New source saved with ID: `%d`\\. Use this ID to manage the source\\.",
				sourceID,
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
