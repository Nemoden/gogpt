package chat

import (
	"context"

	"github.com/nemoden/gogpt/config"
	"github.com/nemoden/gogpt/renderer"
	openai "github.com/sashabaranov/go-openai"
)

type Chat struct {
	config   *config.Config
	messages []openai.ChatCompletionMessage
	client   *openai.Client
	ctx      context.Context
}

func (c *Chat) AddSystemInstruction(m string) {
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: m,
	})
}

func (c *Chat) AskAndRender(input string) error {
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    c.messages,
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
		Stream:      true,
		User:        openai.ChatMessageRoleUser,
	}
	res, err := c.client.CreateChatCompletionStream(c.ctx, req)
	if err != nil {
		return err
	}
	adapter := &renderer.ChatCompletionStreamResponseAdapter{res}
	assistantRespose := c.config.Renderer.Render(adapter)
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: assistantRespose,
	})
	return nil
}

func New(c *config.Config) (*Chat, error) {
	apiKey, err := config.LoadApiKey()
	if err != nil {
		return nil, err
	}
	ch := &Chat{
		config:   c,
		messages: []openai.ChatCompletionMessage{},
		client:   openai.NewClient(apiKey.String()),
		ctx:      context.Background(),
	}
	for _, i := range ch.config.InitialSystemInstructions {
		ch.AddSystemInstruction(i)
	}
	return ch, nil
}
