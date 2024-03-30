package summary

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
	"sync"
)

const (
	aiModel           string  = "gpt-3.5-turbo"
	openAiMaxTokens   int     = 1000
	openAiTemperature float32 = 0.3
	openAiTopP        float32 = 0.6
)

type OpenAIClient struct {
	client  *openai.Client
	enabled bool
	mu      sync.Mutex
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	s := &OpenAIClient{
		client: openai.NewClient(apiKey),
	}

	s.enabled = apiKey != ""
	log.Printf("openai summarizer enabled: %v", s.enabled)

	return s
}

func (s *OpenAIClient) Request(ctx context.Context, text string, prompt string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.enabled {
		return "", nil
	}

	request := openai.ChatCompletionRequest{
		Model: aiModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf("%s%s", text, prompt),
			},
		},
		MaxTokens:   openAiMaxTokens,
		Temperature: openAiTemperature,
		TopP:        openAiTopP,
	}

	resp, err := s.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	rawSummary := strings.TrimSpace(resp.Choices[0].Message.Content)
	return cleanSummary(rawSummary), nil

}

func cleanSummary(rawSummary string) string {
	if strings.HasSuffix(rawSummary, ".") {
		return rawSummary
	}

	sentences := strings.Split(rawSummary, ".")
	return strings.Join(sentences[:len(sentences)-1], ".") + "."
}
