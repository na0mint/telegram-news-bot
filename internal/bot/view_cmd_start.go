package bot

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-bot/internal/botkit"
	"tg-bot/internal/botkit/markup"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error {
		msgText := markup.EscapeForMarkdown("Hello! Use these commands to operate the autoposting bot:" +
			"\n- /sources - get all sources\n- /topics - get all sources" +
			"\n- /addsource {\"name\":\"newSource\",\"url\": \"feed-url\",\"topicID\": \"topic-id\"} - add new source" +
			"\n- /deletesource {sourceId} - delete source by id" +
			"\n- /sourcebyid {sourceId} - get source by id" +
			"\n- /sourcesbytopicid {topicId} - get sources by topic id" +
			"\n- /topics - get all topics",
		)

		reply := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		reply.ParseMode = tgbotapi.ModeMarkdownV2

		if _, err := api.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
