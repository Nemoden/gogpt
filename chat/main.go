package chat

import (
	"github.com/nemoden/gogpt/config"
	openapi "github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openapi.Client
}

func NewClient(c *config.Config) (*Client, error) {
	apiKey, err := config.LoadApiKey()
	if err != nil {
		return nil, err
	}
	client := openapi.NewClient(apiKey.String())
	return &Client{client}, nil
}

func CompletionRequest(prompt string, c *config.Config) openapi.CompletionRequest {
	return openapi.CompletionRequest{
		Model:       c.Model,
		Prompt:      c.PromptPrefix + prompt,
		MaxTokens:   c.MaxTokens,
		Temperature: c.Temperature,
		Stream:      false,
	}
}

func (c *Client) GptClient() *openapi.Client {
	return c.client
}
