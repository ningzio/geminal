package internal

import (
	"fmt"
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var _ Renderer = (*ChromaRenderer)(nil)

func NewChromaRenderer() *ChromaRenderer {
	style := styles.Get("catppuccin-mocha")
	// style := styles.Get("gruvbox")
	// style := styles.Get("github-dark")
	if style == nil {
		style = styles.Fallback
	}
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	return &ChromaRenderer{
		lexer:     lexers.Markdown,
		formatter: formatter,
		style:     style,
	}
}

type ChromaRenderer struct {
	lexer     chroma.Lexer
	formatter chroma.Formatter
	style     *chroma.Style
}

// RenderMessage implements Renderer.
func (cr *ChromaRenderer) RenderMessage(writer io.Writer, message *Message) {
	str := fmt.Sprintf("# %s:\n\n%s\n\n", message.Role, message.Content)

	iterator, err := cr.lexer.Tokenise(nil, str)
	if err != nil {
		writer.Write([]byte(fmt.Sprintf("render: %v", err)))
	}
	if err := cr.formatter.Format(writer, cr.style, iterator); err != nil {
		writer.Write([]byte(fmt.Sprintf("render: %v", err)))
	}
}
