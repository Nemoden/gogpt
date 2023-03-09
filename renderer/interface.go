package renderer

import (
	"bufio"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Choices struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type StreamResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
	Usage   string    `json:"usage"`
}

type StreamResponseAdapter interface {
	Recv() (StreamResponse, error)
	Close()
}

func ReadStreamLineByLine(s StreamResponseAdapter) (string, error) {
	var buffer strings.Builder
	var tok string
	var scanner *bufio.Scanner
	for {
		res, err := s.Recv()
		tok = res.Choices[0].Text
		if len(tok) == 0 {
			continue
		}
		if err != nil {
			return "", err
		}
		buffer.WriteString(tok)
		scanner = bufio.NewScanner(strings.NewReader(buffer.String()))
		for scanner.Scan() {
			return scanner.Text(), nil
		}
	}
}

type ChatCompletionStreamResponseAdapter struct {
	Stream *openai.ChatCompletionStream
}

func (a *ChatCompletionStreamResponseAdapter) Recv() (StreamResponse, error) {
	res, err := a.Stream.Recv()
	if err != nil {
		return StreamResponse{}, err
	}
	choices := make([]Choices, len(res.Choices))
	for i, c := range res.Choices {
		choices[i] = Choices{
			Text:         c.Delta.Content,
			Index:        c.Index,
			FinishReason: c.FinishReason,
		}
	}
	return StreamResponse{
		ID:      res.ID,
		Object:  res.Object,
		Created: res.Created,
		Model:   res.Model,
		Choices: choices,
		Usage:   "",
	}, nil
}

func (a *ChatCompletionStreamResponseAdapter) Close() {
	a.Stream.Close()
}

type CompletionStreamResponseAdapter struct {
	stream *openai.CompletionStream
}

func (a *CompletionStreamResponseAdapter) Recv() (StreamResponse, error) {
	res, err := a.stream.Recv()
	if err != nil {
		return StreamResponse{}, err
	}
	choices := make([]Choices, len(res.Choices))
	for i, c := range res.Choices {
		choices[i] = Choices{
			Text:         c.Text,
			Index:        c.Index,
			FinishReason: c.FinishReason,
		}
	}
	return StreamResponse{
		ID:      res.ID,
		Object:  res.Object,
		Created: res.Created,
		Model:   res.Model,
		Choices: choices,
		Usage:   "",
	}, nil
}

func (a *CompletionStreamResponseAdapter) Close() {
	a.stream.Close()
}

type Renderer interface {
	// Renders response from ChatGPT API to a file. Returns the output as a string.
	Render(stream StreamResponseAdapter) string
}
