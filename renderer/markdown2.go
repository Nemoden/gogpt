package renderer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	markdown "github.com/MichaelMure/go-term-markdown"
)

type Markdown2Renderer struct {
	out    *os.File
	prefix string
}

func NewMarkdown2Renderer(out *os.File, prefix string) *Markdown2Renderer {
	return &Markdown2Renderer{out, prefix}
}

func render(s string) string {
	return string(markdown.Render(s, 80, 4))
}

type rrr struct {
	tokens  []string
	text    string
	curline int
}

func (r *rrr) Add(t string) {
	r.tokens = append(r.tokens, t)
}

func up() {
	fmt.Print("\033[1A\r")
}

func col0() {
	fmt.Print("\033[0G\r")
}

func clearRight() {
	fmt.Print("\033[K\r")
}

func cleanLines(n int) {
	for i := n; i > 0; i-- {
		up()
		time.Sleep(time.Millisecond * 300)
		//fmt.Print("col0")
		col0()
		//fmt.Print("CLEAR:")
		time.Sleep(time.Millisecond * 300)
		clearRight()
		time.Sleep(time.Millisecond * 300)
		//if i != 1 {
		//}
	}
}

func (r *rrr) Render() {
	s := strings.Join(r.tokens, "")
	md := render(s)
	// rendered lines
	mdlines := strings.Split(strings.Trim(md, "\n"), "\n")
	mdlen := len(mdlines)
	if mdlen == 0 {
		return
	}
	r.curline = mdlen
	last := mdlines[mdlen-1]
	//for i := 0; i < mdlen; i++ {
	//fmt.Printf("LineIdx:%d, Line: %s\n", i, mdlines[i])
	//}
	//time.Sleep(600)
	// rendered previous line, start new one
	cleanLines(1)
	//fmt.Printf("%+v %d >>>%s<<< curline %d", mdlines, mdlen, mdlines[mdlen-1], r.curline)
	fmt.Print(last)
	time.Sleep(time.Millisecond * 300)
}

func (r *Markdown2Renderer) Render(stream StreamResponseAdapter) string {
	defer stream.Close()

	//w := ansi.NewWriter(nil)

	wholeResponse := ""
	var token string
	tokens := 0

	rr := &rrr{}

	//wholeBlocksRendered := make([]string, 0)

	curPrintMd := ""
	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			fmt.Printf("\n")
			break
		}

		if len(response.Choices) > 0 {
			tokens += 1
			token = response.Choices[0].Text
			wholeResponse += token

			rr.Add(token)
			rr.Render()
			fmt.Println(curPrintMd)

		}
	}
	return wholeResponse
}

func countLines(b string) int {
	return len(strings.Split(b, "\n"))
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
