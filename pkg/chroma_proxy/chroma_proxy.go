package chroma_proxy

import (
	"bytes"
	"errors"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"golang.org/x/net/html"
)

// Used to strip out chroma prefixing from the CSS.
var reChromaCSS = regexp.MustCompile("(/\\\\*.*\\\\*/)|(\\.chroma)|\\s")

type ProxyHighlighter struct {
	lexer     chroma.Lexer
	style     *chroma.Style
	formatter *chtml.Formatter
}

// NewProxyHighlighter returns a new ProxyHighlighter object that can be used for context aware syntax highlighting.
func NewProxyHighlighter(language, style string) (*ProxyHighlighter, error) {
	p := &ProxyHighlighter{lexer: lexers.Get(language), style: styles.Get(style)}

	if p.lexer == nil {
		return nil, errors.New("could not retrieve chroma lexer for language " + language)
	}
	p.lexer = chroma.Coalesce(p.lexer)

	if p.style == nil {
		p.style = styles.Fallback
	}

	p.formatter = chtml.New(chtml.WithClasses(true), chtml.ClassPrefix(""))

	return p, nil
}

// CSS returns the formatted stylesheet for the chroma highlighting.
func (p *ProxyHighlighter) CSS() (string, error) {
	var buff bytes.Buffer
	if err := p.formatter.WriteCSS(&buff, p.style); err != nil {
		return "", err
	}
	return reChromaCSS.ReplaceAllString(buff.String(), ""), nil
}

// Highlight takes an array of string lines and applies chroma syntax highlighting to them. Context awareness is
// maintained through the highlighting process.
func (p *ProxyHighlighter) Highlight(lines []string) ([]string, error) {
	return p.highlight(lines)
}

func (p *ProxyHighlighter) highlight(lines []string) ([]string, error) {
	file := strings.Join(lines, "\n")

	iterator, err := p.lexer.Tokenise(nil, file)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer
	if err = p.formatter.Format(&result, p.style, iterator); err != nil {
		return nil, err
	}

	return chromaLines(&result)
}

func chromaLines(chromaHTML *bytes.Buffer) ([]string, error) {
	doc, err := html.Parse(chromaHTML)
	if err != nil {
		return nil, err
	}

	code := doc.FirstChild.LastChild.FirstChild.FirstChild.FirstChild
	line, err := chromaLine(code.FirstChild)
	if err != nil {
		return nil, err
	}
	lines := []string{line}

	sibling := code.NextSibling
	for sibling != nil {
		line, err = chromaLine(sibling)
		if err != nil {
			return nil, err
		}
		lines = append(lines, line)
		sibling = sibling.NextSibling
	}

	return lines, nil
}

func chromaLine(node *html.Node) (string, error) {
	var builder strings.Builder
	if err := html.Render(&builder, node); err != nil {
		return "", err
	}
	return strings.ReplaceAll(builder.String(), "\n", ""), nil
}
