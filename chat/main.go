package chat

import (
	"github.com/nemoden/chat/config"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type Client struct {
	client *gogpt.Client
}

func NewClient(c *config.Config) (*Client, error) {
	apiKey, _, err := config.LoadApiKey()
	if err != nil {
		return nil, err
	}
	client := gogpt.NewClient(apiKey)
	return &Client{client}, nil
}

func (c *Client) GptClient() *gogpt.Client {
	return c.client
}
