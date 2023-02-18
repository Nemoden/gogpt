package renderer

import (
	"errors"
	"fmt"
	"io"
	"os"

	gogpt "github.com/sashabaranov/go-gpt3"
)

type TokenRenderer struct {
	out    *os.File
	prefix string
}

// This renderer can be used to view chatgpt response token-by-token
func NewTokenRenderer(out *os.File, prefix string) *TokenRenderer {
	return &TokenRenderer{out, prefix}
}

func (r *TokenRenderer) Render(stream *gogpt.CompletionStream) string {
	var token string
	wholeResponse := ""
	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			fmt.Printf("\n")
			break
		}

		if len(response.Choices) > 0 {
			token = response.Choices[0].Text
			wholeResponse += token
			fmt.Printf("%v -- %s", []byte(token), token)
			fmt.Printf("\n")
		}
	}

	return wholeResponse
}
