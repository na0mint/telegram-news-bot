package bot

import (
	"fmt"
	"tg-bot/internal/botkit/markup"
	"tg-bot/internal/model"
)

func FormatSource(source model.Source) string {
	return fmt.Sprintf(
		"🌐 *%s*\nID: `%d`\nFeed URL: %s\nTopic ID: `%d`",
		markup.EscapeForMarkdown(source.Name),
		source.ID,
		markup.EscapeForMarkdown(source.FeedURL),
		source.TopicID,
	)
}

func FormatTopic(topic model.Topic) string {
	return fmt.Sprintf(
		"💡 *%s*\nID: `%d`\nDescription: %s",
		markup.EscapeForMarkdown(topic.Name),
		topic.ID,
		markup.EscapeForMarkdown(topic.Description),
	)
}
