package renderer

import (
	openai "github.com/sashabaranov/go-openai"
)

type Renderer interface {
	// Renders response from ChatGPT API to a file. Returns the output as a string.
	Render(stream *openai.CompletionStream) string
}
