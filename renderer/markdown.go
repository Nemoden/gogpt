package renderer

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/gosuri/uilive"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type MarkdownRenderer struct {
	out    *os.File
	prefix string
}

func NewMarkdownRenderer(out *os.File, prefix string) *MarkdownRenderer {
	return &MarkdownRenderer{out, prefix}
}

func (r *MarkdownRenderer) Render(stream *gogpt.CompletionStream) string {
	glamourRenderer, _ := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithEmoji(), glamour.WithPreservedNewLines())
	writer := uilive.New()
	writer.Out = r.out
	writer.Start()
	defer writer.Stop()
	previousResponse := ""
	var currentResponse string
	wholeResponse := ""
	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			fmt.Fprintf(writer, "EOF: Printing a new line")
			fmt.Fprintf(writer, "\n")
			previousResponse = ""
			break
		}

		if len(response.Choices) > 0 {
			wholeResponse += response.Choices[0].Text
			currentResponse = previousResponse + response.Choices[0].Text
			out, _ := glamourRenderer.Render(currentResponse)
			fmt.Fprintf(writer, out)
			previousResponse = currentResponse
		}
	}
	return wholeResponse
}
