package renderer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"

	"github.com/nemoden/uilive"
)

type Markdown2Renderer struct {
	out    *os.File
	prefix string
}

func NewMarkdown2Renderer(out *os.File, prefix string) *Markdown2Renderer {
	return &Markdown2Renderer{out, prefix}
}

func (r *Markdown2Renderer) Render(stream StreamResponseAdapter) string {
	defer stream.Close()
	writer := uilive.New()
	writer.Out = r.out
	writer.Start()

	wholeResponse := ""
	currentBlock := ""
	var token string
	var lastToken string
	tokens := 0
	blocks := 0

	//wholeBlocksRendered := make([]string, 0)

	curPrintMd := ""
	insideMarkdownCodeBlockBytes := 0
	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			//wholeBlocksRendered = append(wholeBlocksRendered, currentBlock)
			writer.Stop()
			fmt.Printf("\n")
			break
		}

		if len(response.Choices) > 0 {
			tokens += 1
			token = response.Choices[0].Text
			wholeResponse += token

			if token == "\n" && lastToken == "\n" && insideMarkdownCodeBlockBytes == 0 {
				blocks += 1
				writer.Stop()
				// Start new writer
				writer = uilive.New()
				writer.Out = r.out
				writer.Start()
				//wholeBlocksRendered = append(wholeBlocksRendered, currentBlock)
				currentBlock = ""
				continue
			}
			currentBlock += token
			if insideMarkdownCodeBlockBytes == 0 {
				insideMarkdownCodeBlockBytes = isInsideMarkdownCodeBlock(currentBlock)
			}
			if insideMarkdownCodeBlockBytes > 0 && strings.HasSuffix(currentBlock, "```") && len(currentBlock) > (insideMarkdownCodeBlockBytes+3) {
				insideMarkdownCodeBlockBytes = 0
			}

			if insideMarkdownCodeBlockBytes > 0 {
				curPrintMd = string(markdown.Render(currentBlock+"\n```", 80, 4))
			} else {
				curPrintMd = string(markdown.Render(currentBlock, 80, 4))
			}
			fmt.Fprintln(writer, curPrintMd)
			//fmt.Fprintln(writer, currentBlock)

			lastToken = token
		}
	}
	//for _, b := range wholeBlocksRendered {
	//fmt.Printf("---\n%s\n---\n", b)
	//}
	return wholeResponse
}

func isInsideMarkdownCodeBlock(str string) int {
	r := strings.NewReader(str)
	s := bufio.NewScanner(r)
	hasTicks := func(s string) bool {
		return strings.HasPrefix(s, "```")
	}
	lines := 2
	bytesRead := 0
	var text string
	for i := 0; i < lines; i += 1 {
		if s.Scan() {
			text = s.Text()
			bytesRead += len(text) + 1 // \n
			if hasTicks(s.Text()) {
				return bytesRead
			}
		}
	}
	return 0
}
