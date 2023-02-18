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

func NewTokenRenderer(out *os.File, prefix string) *TokenRenderer {
	return &TokenRenderer{out, prefix}
}

func (r *TokenRenderer) Render(stream *gogpt.CompletionStream) string {
	var token string
	wholeResponse := ""
	var lastToken string
	cutoff := false
	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			fmt.Printf("EOF: Printing a new line")
			fmt.Printf("\n")
			break
		}

		if len(response.Choices) > 0 {
			token = response.Choices[0].Text
			wholeResponse += token
			fmt.Printf("%v --- %s --- is NL %v --- is cutoff %v", []byte(token), token, token == "\n", cutoff)
			fmt.Printf("\n")

			if token == "\n" && lastToken == "\n" {
				cutoff = true
			} else {
				cutoff = false
			}
			lastToken = token
		}
	}

	fmt.Printf("Whole response:\n%s\n", wholeResponse)
	return wholeResponse
}