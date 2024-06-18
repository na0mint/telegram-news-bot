package summary

import (
	"context"
	"fmt"
	"log"
	"sync"
	"tg-bot/internal/config"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const (
	llmModel = "llama3"
)

type LocalLLM struct {
	llm     *ollama.LLM
	enabled bool
	mu      sync.Mutex
}

func NewLocalLLM() (*LocalLLM, error) {
	llm, err := ollama.New(ollama.WithModel(llmModel))
	if err != nil {
		return nil, err
	}

	s := &LocalLLM{
		llm:     llm,
		enabled: config.Get().IsLocalLLM,
	}

	log.Printf("local llm enabled: %v, model: %v", s.enabled, llmModel)
	return s, nil
}

func (l *LocalLLM) Request(ctx context.Context, text string, prompt string) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.enabled {
		return "", nil
	}

	query := fmt.Sprintf("%s%s", prompt, text)

	completion, err := llms.GenerateFromSinglePrompt(ctx, l.llm, query)
	if err != nil {
		log.Printf("[ERROR] failed to generate summary: %v", err)
		return "", err
	}

	return completion, nil
}
