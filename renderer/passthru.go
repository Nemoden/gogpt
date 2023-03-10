package renderer

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// Passthru renderer simlpy passes thru the response from ChatGPT API and renders it as is.
type PassthruRenderer struct {
	out    *os.File
	prefix string
}

func NewPassthruRenderer(out *os.File, prefix string) *PassthruRenderer {
	return &PassthruRenderer{out, prefix}
}

func (r *PassthruRenderer) Render(stream StreamResponseAdapter) string {
	defer stream.Close()
	wholeResponse := ""
	writer := bufio.NewWriter(r.out)
	for {
		response, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			writer.WriteString("\n")
			writer.Flush()
			break
		}

		if len(response.Choices) > 0 {
			wholeResponse += response.Choices[0].Text
			writer.WriteString(response.Choices[0].Text)
			writer.Flush()
		}
	}
	return wholeResponse
}
