package chat

import (
	"github.com/nemoden/gogpt/config"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type Client struct {
	client *gogpt.Client
}

func NewClient(c *config.Config) (*Client, error) {
	apiKey, err := config.LoadApiKey()
	if err != nil {
		return nil, err
	}
	client := gogpt.NewClient(apiKey.String())
	return &Client{client}, nil
}

func CompletionRequest(prompt string, c *config.Config) gogpt.CompletionRequest {
	return gogpt.CompletionRequest{
		Model:       c.Model,
		Prompt:      c.PromptPrefix + prompt,
		MaxTokens:   c.MaxTokens,
		Temperature: c.Temperature,
		Stream:      false,
	}
}

func (c *Client) GptClient() *gogpt.Client {
	return c.client
}
