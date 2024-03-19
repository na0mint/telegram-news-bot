package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"strconv"
	"strings"
	"tg-bot/internal/botkit"
	"tg-bot/internal/model"
)

type SourceByTopicLister interface {
	SourcesByTopicId(ctx context.Context, topicId int64) ([]model.Source, error)
}

func ViewCmdListSourcesByTopicId(lister SourceByTopicLister) botkit.ViewFunc {
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

		sources, err := lister.SourcesByTopicId(ctx, targetId)
		if err != nil {
			return err
		}

		var (
			sourceInfo = lo.Map(sources, func(source model.Source, _ int) string {
				return FormatSource(source)
			})

			msgText = fmt.Sprintf(
				"Sources \\(total %d\\):\n\n%s",
				len(sources),
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
