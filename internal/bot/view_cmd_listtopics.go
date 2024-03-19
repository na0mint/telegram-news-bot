package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"strings"
	"tg-bot/internal/botkit"
	"tg-bot/internal/model"
)

type TopicLister interface {
	Topics(ctx context.Context) ([]model.Topic, error)
}

func ViewCmdListTopics(lister TopicLister) botkit.ViewFunc {
	return func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error {
		topics, err := lister.Topics(ctx)
		if err != nil {
			return err
		}

		var (
			sourceInfo = lo.Map(topics, func(topic model.Topic, _ int) string {
				return FormatTopic(topic)
			})

			msgText = fmt.Sprintf(
				"Sources \\(total %d\\):\n\n%s",
				len(topics),
				strings.Join(sourceInfo, "\n\n"),
			)
		)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		reply.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := api.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
