package renderer

import (
	gogpt "github.com/sashabaranov/go-gpt3"
)

type Renderer interface {
	// Renders response from ChatGPT API to a file. Returns the output as a string.
	Render(stream *gogpt.CompletionStream) string
}
