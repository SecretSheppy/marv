package server

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/SecretSheppy/marv/pkg/highlighter"
	"github.com/SecretSheppy/marv/pkg/mutations"
	"github.com/alecthomas/chroma/v2/styles"
)

type Renderer struct {
	ext       string
	lines     []string
	conflicts mutations.Conflicts
	highlight *highlighter.Highlighter
	meta      *Meta
}

func NewRenderer(ext string, lines []string, meta *Meta, conflicts mutations.Conflicts) *Renderer {
	return &Renderer{ext: ext, lines: lines, meta: meta, conflicts: conflicts}
}

// Render writes the HTML page into a writer.
func (r *Renderer) Render(w io.Writer) error {
	render, err := r.render()
	if err != nil {
		return err
	}
	if _, err := w.Write(render); err != nil {
		return err
	}
	return nil
}

func (r *Renderer) render() ([]byte, error) {
	var (
		buff bytes.Buffer
		err  error
	)

	r.highlight, err = highlighter.NewHighlighter(r.ext, r.lines, styles.Get("monokai"))
	if err != nil {
		return nil, err
	}

	r.scripts(&buff)
	if err = r.styles(&buff); err != nil {
		return nil, err
	}
	if err = r.code(&buff); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (r *Renderer) scripts(w *bytes.Buffer) {
	for _, script := range r.meta.scripts {
		w.WriteString("<script type=\"text/javascript\">")
		w.Write(script)
		w.WriteString("</script>")
	}
}

func (r *Renderer) styles(w *bytes.Buffer) error {
	w.WriteString("<style>")
	w.Write(r.meta.style)

	code, err := r.highlight.CSS()
	if err != nil {
		return err
	}
	w.WriteString(code)

	w.WriteString("</style>")
	return nil
}

type renderedConflict struct {
	start, end int
	render     string
}

type LineDiffType int

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

const (
	LineRemoved LineDiffType = iota
	LineEqual
	LineInserted
)

func (r *Renderer) codeline(w *bytes.Buffer, ln int, lt LineDiffType, code string) {
	w.WriteString("<tr>")
	w.WriteString(fmt.Sprintf("<td class=\"ln\">%d</td>", ln))
	w.WriteString(fmt.Sprintf("<td class=\"lt\">%s</td>", lt))
	w.WriteString(fmt.Sprintf("<td class=\"code\">%s</td>", code))
	w.WriteString("</tr>")
}

func (r *Renderer) code(w *bytes.Buffer) error {
	r.conflicts.Sort()
	rendered := make([]*renderedConflict, 0, len(r.conflicts))
	for _, conflict := range r.conflicts {
		render, err := r.conflict(conflict)
		if err != nil {
			return err
		}
		rendered = append(rendered, render)
	}

	w.WriteString("<table class=\"bg code\">")
	for i := 0; i < len(r.lines); i++ {
		if len(rendered) > 0 {
			if conflict := rendered[0]; conflict.start == i {
				w.WriteString(conflict.render)
				rendered = rendered[1:]
				i = conflict.end
				continue
			}
		}
		line, err := r.highlight.HighlightLine(i)
		if err != nil {
			return err
		}
		r.codeline(w, i+1, LineEqual, line)
	}
	w.WriteString("</table>")
	return nil
}

func (r *Renderer) conflict(c *mutations.Conflict) (*renderedConflict, error) {
	var builder strings.Builder
	for _, mutation := range c.Mutations {
		render, err := r.mutation(c.StartLine, c.EndLine, mutation)
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

func (r *Renderer) mutation(start, end int, m *mutations.Mutation) (string, error) {
	var buff bytes.Buffer
	for i := m.Starts.Line; i <= m.Ends.Line; i++ {
		var pre, post string
		diff := r.lines[i]
		if i == m.Ends.Line {
			post = diff[m.Ends.Char:]
			diff = diff[:m.Ends.Char]
		}
		if i == m.Starts.Line {
			pre = diff[:m.Starts.Char]
			diff = diff[m.Starts.Char:]
		}
		hl, err := highlighter.NewHighlighter(r.highlight.Lang(), []string{pre, diff, post}, r.highlight.Style())
		if err != nil {
			return "", err
		}
		lines, err := hl.HighlightLines(0, 2)
		if err != nil {
			return "", err
		}
		for j := 0; j < len(lines); j++ {
			if line := lines[j]; line != "" {
				lines[j] = line[19 : len(line)-7]
			}
		}
		code := fmt.Sprintf("%s<span class=\"highlight\">%s</span>%s", lines[0], lines[1], lines[2])
		r.codeline(&buff, i+1, LineRemoved, code)
	}

	// TODO: highlight inserted lines

	// TODO: pad (if necessary) with lines from before and after the mutation so that it is the same lines as the
	//  conflict.

	// FIXME: all mutant blocks require same size to make sense
	//str := fmt.Sprintf("<tr><td>desc: %s, mutant: \"%s\"</td></tr>", m.Name, m.Source)
	return buff.String(), nil
}
