package botkit

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"runtime/debug"
	"time"
)

const updateTimeout int = 60

type Bot struct {
	api      *tgbotapi.BotAPI
	cmdViews map[string]ViewFunc
}

func NewBot(api *tgbotapi.BotAPI) *Bot {
	return &Bot{api: api}
}

type ViewFunc func(ctx context.Context, api *tgbotapi.BotAPI, update tgbotapi.Update) error

func (b *Bot) RegisterCmdView(cmd string, view ViewFunc) {
	if b.cmdViews == nil {
		b.cmdViews = make(map[string]ViewFunc)
	}

	b.cmdViews[cmd] = view
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = updateTimeout

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, 5*time.Second)
			b.handleUpdate(updateCtx, update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) SendMessage(chatId int64, message string) {
	if _, err := b.api.Send(
		tgbotapi.NewMessage(chatId, message),
	); err != nil {
		log.Printf("[ERROR] failed to send message: %v to chat: %v", message, chatId)
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recoverd: %v\n%s", p, string(debug.Stack()))
		}
	}()

	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	if !update.Message.IsCommand() {
		b.SendMessage(update.Message.Chat.ID, "Your message is not a command")
		return
	}

	cmd := update.Message.Command()

	cmdView, ok := b.cmdViews[cmd]
	if !ok {
		b.SendMessage(update.Message.Chat.ID, "Unable to find your command")
		return
	}

	if err := cmdView(ctx, b.api, update); err != nil {
		log.Printf("[ERROR] failed to handle update: %v", err)
		b.SendMessage(update.Message.Chat.ID, "internal error")
	}
}
