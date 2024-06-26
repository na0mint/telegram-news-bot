package notifier

import (
	"context"
	"fmt"
	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"
	"io"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"tg-bot/internal/botkit/markup"
	"tg-bot/internal/config"
	"tg-bot/internal/model"
	"time"
)

var (
	NewLinesRegexp = regexp.MustCompile(`\n{3,}`)
)

const (
	articlesOffset int64  = 1000
	translation    string = "translation"
)

type ArticleProvider interface {
	FindAllNotPosted(ctx context.Context, since time.Time, limit int64) ([]model.Article, error)
	MarkPostedById(ctx context.Context, id int64) error
}

type SourceProvider interface {
	SourcesByTopicId(ctx context.Context, topicId int64) ([]model.Source, error)
	Sources(ctx context.Context) ([]model.Source, error)
}

type AIClient interface {
	Request(ctx context.Context, text string, prompt string) (string, error)
}

type Notifier struct {
	articles         ArticleProvider
	sources          SourceProvider
	openAIClient     AIClient
	bot              *tgbotapi.BotAPI
	sendInterval     time.Duration
	lookupTimeWindow time.Duration
	channelId        int64
}

func NewNotifier(
	articles ArticleProvider,
	sources SourceProvider,
	summarizer AIClient,
	bot *tgbotapi.BotAPI,
	sendInterval time.Duration,
	lookupTimeWindow time.Duration,
	channelId int64,
) *Notifier {
	return &Notifier{
		articles:         articles,
		sources:          sources,
		openAIClient:     summarizer,
		bot:              bot,
		sendInterval:     sendInterval,
		lookupTimeWindow: lookupTimeWindow,
		channelId:        channelId,
	}
}

func (n *Notifier) Start(ctx context.Context) error {
	ticker := time.NewTicker(n.sendInterval)
	defer ticker.Stop()

	if err := n.SelectAndSendArticle(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := n.SelectAndSendArticle(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topArticles, err := n.articles.FindAllNotPosted(ctx, time.Now().Add(-n.lookupTimeWindow), articlesOffset)
	if err != nil {
		return err
	}

	if len(topArticles) == 0 {
		return nil
	}

	sources, err := n.sources.Sources(ctx)
	if err != nil {
		return err
	}

	for _, topicId := range getUniqueTopicIds(sources) {

		sourcesForTopicId := lo.Filter(sources, func(source model.Source, _ int) bool {
			return source.TopicID == topicId
		})

		sourceIds := lo.Map(sourcesForTopicId, func(source model.Source, _ int) int64 {
			return source.ID
		})

		article := lo.Filter(topArticles, func(article model.Article, _ int) bool {
			return slices.Contains(sourceIds, article.SourceID)
		})[0]

		postSource, _ := lo.Find(sourcesForTopicId, func(source model.Source) bool {
			return source.ID == article.SourceID
		})

		postText, err := n.extractSummary(ctx, article, postSource.Type)
		if err != nil {
			return err
		}

		if err := n.sendArticle(postText, article); err != nil {
			return err
		}

		if err := n.articles.MarkPostedById(ctx, article.ID); err != nil {
			return err
		}
	}

	return nil
}

func getUniqueTopicIds(sources []model.Source) []int64 {
	topicIds := make([]int64, 0, len(sources))

	for _, source := range sources {
		topicIds = append(topicIds, source.TopicID)
	}

	return lo.Uniq(topicIds)
}

func (n *Notifier) extractSummary(ctx context.Context, article model.Article, postType string) (string, error) {
	var reader io.Reader

	if article.Summary != "" {
		reader = strings.NewReader(article.Summary)
	} else {
		response, err := http.Get(article.Link)
		if err != nil {
			return "", err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("[ERROR] Failed to close response body: %v", err)
			}
		}(response.Body)

		reader = response.Body
	}

	doc, err := readability.FromReader(reader, nil)
	if err != nil {
		return "", err
	}

	return n.makeSummary(ctx, doc, postType)
}

func (n *Notifier) makeSummary(ctx context.Context, article readability.Article, postType string) (string, error) {
	switch postType {
	case translation:
		translation, err := n.openAIClient.Request(ctx, cleanText(article.TextContent), config.Get().AITranslationPrompt)
		if err != nil {
			return "", err
		}

		log.Println("Translation post")

		return "\n\n" + translation, nil
	default:
		summary, err := n.openAIClient.Request(ctx, cleanText(article.TextContent), config.Get().AIDefaultPrompt)
		if err != nil {
			return "", err
		}

		log.Println("Summary post")

		return "\n\n" + summary, nil
	}
}

func cleanText(text string) string {
	return NewLinesRegexp.ReplaceAllString(text, "\n")
}

func (n *Notifier) sendArticle(summary string, article model.Article) error {
	const msgFormat = "*%s*%s\n\n%s"

	msg := tgbotapi.NewMessage(n.channelId, fmt.Sprintf(
		msgFormat,
		markup.EscapeForMarkdown(article.Title),
		markup.EscapeForMarkdown(summary),
		markup.EscapeForMarkdown(article.Link),
	))
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := n.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
