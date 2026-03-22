package highlighter

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

var ErrNoLexer = errors.New("no lexer for supplied language")

// Highlighter is a wrapper for the chroma library that allows for line by line and substring highlighting.
type Highlighter struct {
	lang      string
	lines     []string
	lexer     chroma.Lexer
	style     *chroma.Style
	formatter *html.Formatter
}

// NewHighlighter returns a newly configured instance of Highlighter.
func NewHighlighter(lang string, lines []string, style *chroma.Style) (*Highlighter, error) {
	h := &Highlighter{
		lang:      lang,
		lines:     lines,
		lexer:     lexers.Get(lang),
		style:     style,
		formatter: html.New(html.WithClasses(true), html.ClassPrefix("")),
	}

	if h.lexer == nil {
		return nil, ErrNoLexer
	}
	h.lexer = chroma.Coalesce(h.lexer)

	if h.style == nil {
		h.style = styles.Fallback
	}

	return h, nil
}

// CSS returns the formatted CSS classes for the chroma highlighting without being wrapped in <style> tags.
func (h *Highlighter) CSS() (string, error) {
	var result bytes.Buffer
	if err := h.formatter.WriteCSS(&result, h.style); err != nil {
		return "", err
	}
	return result.String(), nil
}

func (h *Highlighter) highlightStr(line string) (string, error) {
	iterator, err := h.lexer.Tokenise(nil, line)
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	if err := h.formatter.Format(&result, h.style, iterator); err != nil {
		return "", err
	}
	highlighted := result.String()
	return highlighted[26 : len(highlighted)-13], nil
}

// HighlightLine highlights the specified line number from the sources lines that were provided to the Highlighter
// instance.
func (h *Highlighter) HighlightLine(line int) (string, error) {
	return h.highlightStr(h.lines[line])
}

// HighlightLines calls HighlightLine over a range of lines and returns a list of the results.
func (h *Highlighter) HighlightLines(start, end int) ([]string, error) {
	var err error
	lines := make([]string, end-start)
	for line := start; line < end; line++ {
		lines[line], err = h.HighlightLine(line)
		if err != nil {
			return nil, err
		}
	}
	return lines, nil
}

func (h *Highlighter) highlightSubstr(line string, start, end int, context bool) (string, error) {
	if !context {
		return h.highlightStr(line[start:end])
	}
	prefix, err := h.highlightStr(line[:start])
	if err != nil {
		return "", err
	}
	result, err := h.highlightStr(line[:end])
	if err != nil {
		return "", err
	}
	return result[len(prefix):], nil
}

// HighlightSubstr highlights and returns the specified substring of a given line. If the context value is set to true,
// the Highlighter will highlight the whole line and then identify and return the specified substring HTML elements.
func (h *Highlighter) HighlightSubstr(line, start, end int, context bool) (string, error) {
	return h.highlightSubstr(h.lines[line], start, end, context)
}

func (h *Highlighter) replaceSubstr(replacement string, line, start, end int) string {
	return fmt.Sprintf("%s%s%s", h.lines[line][:start], replacement, h.lines[line][end:])
}

// HighlightReplacedSubstr replaces the specified substring in the line with the newly provided string. It then
// highlights and returns the specified substring of a given line. If the context value is set to true, the Highlighter
// will highlight the whole line and then identify and return the specified substring HTML elements.
func (h *Highlighter) HighlightReplacedSubstr(replacement string, line, start, end int, context bool) (string, error) {
	return h.highlightSubstr(h.replaceSubstr(replacement, line, start, end), start, start+len(replacement), context)
}
