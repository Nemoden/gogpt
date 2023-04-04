package renderer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

type MarkdownRenderer struct {
	out    *os.File
	prefix string
}

func NewMarkdownRenderer(out *os.File, prefix string) *MarkdownRenderer {
	return &MarkdownRenderer{out, prefix}
}

func LineByLine(stream StreamResponseAdapter) <-chan string {
	ch := make(chan string)
	go (func() {
		defer stream.Close()
		glamourRenderer, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithEmoji(),
			glamour.WithPreservedNewLines(),
		)
		var token string
		whole := ""
		var md string
		curpos := 0
		lcount := 0
		wholepos := 0
		insideCode := false
		var torender string
		for {
			res, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Printf("%s", md)
				//fmt.Printf("LCOUNT: %d\n", lcount)
				close(ch)
				return
			}
			if len(res.Choices) == 0 {
				continue
			}
			token = res.Choices[0].Text
			whole += token

			wholeline, err := readLine(whole, &wholepos)
			if err == nil {
				if wholeline == "```\n" {
					if insideCode {
						//fmt.Printf("OUTSIDE THE CODE!")
						insideCode = false
					} else {
						//fmt.Printf("INSIDE THE CODE!")
						insideCode = true
						continue
					}
				}
			}
			if insideCode && err != nil {
				// skip ahead to read a whole line of code
				continue
			}
			torender = whole
			if insideCode {
				torender += "\n```\n"
			}
			md, _ = glamourRenderer.Render(torender)

			if !(insideCode || wholeline == "```\n") {
				md = strings.TrimRight(md, "\n")
			}

			mdline, err := readLine(md, &curpos)
			if err == nil {
				lcount += 1
				ch <- mdline
			}
		}
	})()
	return ch
}

func readLine(str string, pos *int) (string, error) {
	buf := ""
	for i := *pos; i < len(str); i++ {
		buf += string(str[i])
		if str[i] == '\n' {
			// next pos.
			*pos = i + 1
			return buf, nil
		}
	}
	return "", errors.New("No EOF")
}

func (r *MarkdownRenderer) Render(stream StreamResponseAdapter) string {
	for l := range LineByLine(stream) {
		fmt.Printf("%s", l)
		//fmt.Println("--------")
		//fmt.Printf("%v", []byte(l))
		//fmt.Println("")
		//fmt.Printf("%s", l)
		//fmt.Println("\n--------")
	}
	return "sasasasas"
}
