package html

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/pkg/highlighter"
	"github.com/alecthomas/chroma/v2/styles"
)

type renderedConflict struct {
	start, end int
	render     string
}

type LineDiffType int

const (
	LineRemoved LineDiffType = iota
	LineEqual
	LineInserted
)

func (l LineDiffType) String() string {
	switch l {
	case LineRemoved:
		return "-"
	case LineInserted:
		return "+"
	default:
		return ""
	}
}

func (l LineDiffType) CSSClass() string {
	switch l {
	case LineRemoved:
		return "remove"
	case LineInserted:
		return "insert"
	default:
		return ""
	}
}

type CodeRenderer struct {
	ext       string
	lines     []string
	conflicts mutations.Conflicts
	highlight *highlighter.Highlighter
}

func NewCodeRenderer(ext string, lines []string, conflicts mutations.Conflicts) (*CodeRenderer, error) {
	r := &CodeRenderer{ext: ext, lines: lines, conflicts: conflicts}
	var err error
	r.highlight, err = highlighter.NewHighlighter(r.ext, r.lines, styles.Get("darcula"))
	return r, err
}

func (r *CodeRenderer) SyntaxHighlighter() *highlighter.Highlighter {
	return r.highlight
}

func (r *CodeRenderer) Render(w *bytes.Buffer) error {
	render, err := r.render()
	if err != nil {
		return err
	}
	if _, err := w.Write(render); err != nil {
		return err
	}
	return nil
}

func (r *CodeRenderer) render() ([]byte, error) {
	r.conflicts.Sort()
	rendered := make([]*renderedConflict, 0, len(r.conflicts))
	for _, conflict := range r.conflicts {
		render, err := r.renderConflict(conflict)
		if err != nil {
			return nil, err
		}
		rendered = append(rendered, render)
	}

	var buff bytes.Buffer
	buff.WriteString("<table class=\"bg code\">")
	for i := 0; i < len(r.lines); i++ {
		if len(rendered) > 0 {
			if conflict := rendered[0]; conflict.start == i {
				buff.WriteString(conflict.render)
				rendered = rendered[1:]
				i = conflict.end
				continue
			}
		}
		line, err := r.highlight.HighlightLine(i)
		if err != nil {
			return nil, err
		}
		r.renderLine(&buff, i+1, LineEqual, line)
	}
	buff.WriteString("</table>")
	return buff.Bytes(), nil
}

func (r *CodeRenderer) renderLine(w *bytes.Buffer, ln int, lt LineDiffType, code string) {
	w.WriteString(fmt.Sprintf("<tr class=\"%s\">", lt.CSSClass()))
	w.WriteString("<td class=\"line-number\">")
	if ln != 0 {
		w.WriteString(fmt.Sprintf("%d", ln))
	}
	w.WriteString("</td>")
	w.WriteString(fmt.Sprintf("<td class=\"line-type\">%s</td>", lt))
	w.WriteString(fmt.Sprintf("<td class=\"line-content\">%s</td>", code))
	w.WriteString("</tr>")
}

func (r *CodeRenderer) renderConflict(c *mutations.Conflict) (*renderedConflict, error) {
	var builder strings.Builder
	for _, mutation := range c.Mutations {
		render, err := r.renderMutation(c, mutation)
		if err != nil {
			return nil, err
		}
		builder.WriteString(render)
	}
	return &renderedConflict{
		start:  c.StartLine,
		end:    c.EndLine,
		render: builder.String(),
	}, nil
}

func (r *CodeRenderer) renderMutation(c *mutations.Conflict, m *mutations.Mutation) (string, error) {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("<tbody id=\"%s\" data-conflict-id=\"%s\" data-status=\"%s\" data-class=\"mutant\">", m.ID, c.ID, m.Status.Text()))
	r.renderMutationHeader(&buff, m)

	for i := c.StartLine; i < m.Start.Line; i++ {
		line, err := r.highlight.HighlightLine(i)
		if err != nil {
			return "", err
		}
		r.renderLine(&buff, i+1, LineEqual, line)
	}

	for i := m.Start.Line; i <= m.End.Line; i++ {
		var pre, post string
		diff := r.lines[i]
		if i == m.End.Line {
			post = diff[m.End.Char:]
			diff = diff[:m.End.Char]
		}
		if i == m.Start.Line {
			pre = diff[:m.Start.Char]
			diff = diff[m.Start.Char:]
		}
		lines, err := r.highlightMutationParts(pre, diff, post)
		if err != nil {
			return "", err
		}
		code := fmt.Sprintf("%s<span class=\"highlight remove\">%s</span>%s", lines[0], lines[1], lines[2])
		r.renderLine(&buff, i+1, LineRemoved, code)
	}

	mutLines := make([]string, 0)
	for line := range strings.Lines(m.Replacement) {
		mutLines = append(mutLines, strings.ReplaceAll(line, "\n", ""))
	}
	for i, diff := range mutLines {
		var pre, post string
		if i == len(mutLines)-1 { // NOTE: last mutated line
			post = r.lines[m.End.Line][m.End.Char:]
		}
		if i == 0 { // NOTE: first mutated line
			pre = r.lines[m.Start.Line][:m.Start.Char]
		}
		lines, err := r.highlightMutationParts(pre, diff, post)
		if err != nil {
			return "", err
		}
		code := fmt.Sprintf("%s<span class=\"highlight insert\">%s</span>%s", lines[0], lines[1], lines[2])
		r.renderLine(&buff, 0, LineInserted, code)
	}

	for i := m.End.Line; i < c.EndLine; i++ {
		line, err := r.highlight.HighlightLine(i)
		if err != nil {
			return "", err
		}
		r.renderLine(&buff, i+1, LineEqual, line)
	}

	buff.WriteString("</tbody>")
	return buff.String(), nil
}

func (r *CodeRenderer) highlightMutationParts(pre, diff, post string) ([]string, error) {
	hl, err := highlighter.NewHighlighter(r.highlight.Lang(), []string{pre, diff, post}, r.highlight.Style())
	if err != nil {
		return nil, err
	}
	lines, err := hl.HighlightLines(0, 2)
	if err != nil {
		return nil, err
	}
	for j := 0; j < len(lines); j++ {
		if line := lines[j]; line != "" {
			lines[j] = line[19 : len(line)-7]
		}
	}
	return lines, nil
}

func (r *CodeRenderer) renderMutationHeader(buff *bytes.Buffer, m *mutations.Mutation) {
	// TODO: round outside corners and inside corners and give it a border etc...
	buff.WriteString("<tr><td colspan=\"100%\"><div class=\"mutation-header\">")
	buff.WriteString(m.Status.IconWithText())
	buff.WriteString(html.EscapeString(m.Description))
	buff.WriteString("</div></td></tr>")
}
