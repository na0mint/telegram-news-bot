package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/internal/botkit"
	"tg-bot/internal/model"
)

type TopicStorage interface {
	Save(ctx context.Context, topic model.Topic) (int64, error)
}

func ViewCmdAddTopic(storage TopicStorage) botkit.ViewFunc {
	type addTopicArgs struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	return func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args, err := botkit.ParseJSON[addTopicArgs](update.Message.CommandArguments())
		if err != nil {
			_, sendErr := api.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
				"Failed to parse command arguments"))
			if sendErr != nil {
				return sendErr
			}

			return err
		}

		topic := model.Topic{
			Name:        args.Name,
			Description: args.Description,
		}

		topicID, err := storage.Save(ctx, topic)
		if err != nil {
			return err
		}

		var (
			msgText = fmt.Sprintf(
				"New topic saved with ID: `%d`\\. Use this ID to manage the topic\\.",
				topicID,
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
