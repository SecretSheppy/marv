package html

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/internal/review"
	"github.com/SecretSheppy/marv/pkg/highlighter"
	"github.com/alecthomas/chroma/v2/styles"
	"gorm.io/gorm"
)

type renderedConflict struct {
	start, end int
	render     string
}

type lineDiffType int

const (
	lineRemoved lineDiffType = iota
	lineEqual
	lineInserted
)

func (l lineDiffType) String() string {
	switch l {
	case lineRemoved:
		return "-"
	case lineInserted:
		return "+"
	default:
		return ""
	}
}

func (l lineDiffType) CSSClass() string {
	switch l {
	case lineRemoved:
		return "remove"
	case lineInserted:
		return "insert"
	default:
		return ""
	}
}

type codeRendererConfig struct {
	RenderAllMutantData bool
}

type codeRenderer struct {
	ext, framework, file string
	lines                []string
	conflicts            mutations.Conflicts
	highlight            *highlighter.Highlighter
	lnPadding            int
	config               *codeRendererConfig
	db                   *review.Repository
}

func newCodeRenderer(ext, framework, file string, lines []string, conflicts mutations.Conflicts, config *codeRendererConfig, db *review.Repository) (*codeRenderer, error) {
	r := &codeRenderer{
		ext:       ext,
		framework: framework,
		file:      file,
		lines:     lines,
		conflicts: conflicts,
		config:    config,
		db:        db,
	}
	var err error
	r.highlight, err = highlighter.NewHighlighter(r.ext, r.lines, styles.Get("darcula"))
	return r, err
}

func (r *codeRenderer) SyntaxHighlighter() *highlighter.Highlighter {
	return r.highlight
}

func (r *codeRenderer) Render(w *bytes.Buffer) error {
	render, err := r.render()
	if err != nil {
		return err
	}
	if _, err := w.Write(render); err != nil {
		return err
	}
	return nil
}

func (r *codeRenderer) render() ([]byte, error) {
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
	buff.WriteString("<table id=\"code-table\" class=\"bg code\">")
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
		r.renderLine(&buff, i+1, lineEqual, line)
	}
	buff.WriteString("</table>")
	return buff.Bytes(), nil
}

func (r *codeRenderer) renderLine(w *bytes.Buffer, lineNumber int, lt lineDiffType, code string) {
	w.WriteString(fmt.Sprintf("<tr class=\"%s\">", lt.CSSClass()))
	w.WriteString("<td class=\"line-number\">")
	if lineNumber != 0 {
		// NOTE: Padding is used to ensure the line number column is the same width through all the individual <tbody>
		// elements.
		w.WriteString(fmt.Sprintf("%*d", r.padding(), lineNumber))
	}
	w.WriteString("</td>")
	w.WriteString(fmt.Sprintf("<td class=\"line-type\">%s</td>", lt))
	w.WriteString(fmt.Sprintf("<td class=\"line-content\">%s</td>", code))
	w.WriteString("</tr>")
}

func (r *codeRenderer) padding() int {
	if r.lnPadding == 0 {
		r.lnPadding = len(strconv.Itoa(len(r.lines)))
	}
	return r.lnPadding
}

func (r *codeRenderer) renderConflict(c *mutations.Conflict) (*renderedConflict, error) {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("<tbody class=\"hidden\" data-conflict-id=\"%s\">", c.ID))
	for i := c.StartLine; i <= c.EndLine; i++ {
		line, err := r.highlight.HighlightLine(i)
		if err != nil {
			return nil, err
		}
		r.renderLine(&buff, i+1, lineEqual, line)
	}
	buff.WriteString("</tbody>")
	for _, mutation := range c.Mutations {
		render, err := r.renderMutation(c, mutation)
		if err != nil {
			return nil, err
		}
		buff.WriteString(render)
	}
	return &renderedConflict{
		start:  c.StartLine,
		end:    c.EndLine,
		render: buff.String(),
	}, nil
}

func (r *codeRenderer) renderMutation(c *mutations.Conflict, m *mutations.Mutation) (string, error) {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("<tbody id=\"%s\" data-conflict-id=\"%s\" data-status=\"%s\" data-class=\"mutant\" class=\"mutation\">", m.ID, c.ID, m.Status.Text()))
	r.renderMutationHeader(&buff, m)
	if r.config.RenderAllMutantData {
		r.renderAllMutationData(&buff, m)
	}

	for i := c.StartLine; i < m.Start.Line; i++ {
		line, err := r.highlight.HighlightLine(i)
		if err != nil {
			return "", err
		}
		r.renderLine(&buff, i+1, lineEqual, line)
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
		r.renderLine(&buff, i+1, lineRemoved, code)
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
		r.renderLine(&buff, 0, lineInserted, code)
	}

	for i := m.End.Line + 1; i <= c.EndLine; i++ {
		line, err := r.highlight.HighlightLine(i)
		if err != nil {
			return "", err
		}
		r.renderLine(&buff, i+1, lineEqual, line)
	}

	if err := r.renderReviewField(&buff, m); err != nil {
		return "", err
	}

	buff.WriteString("</tbody>")
	return buff.String(), nil
}

func (r *codeRenderer) highlightMutationParts(pre, diff, post string) ([]string, error) {
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

func (r *codeRenderer) renderMutationHeader(buff *bytes.Buffer, m *mutations.Mutation) {
	buff.WriteString("<tr><td colspan=\"100%\"><div class=\"mutation-header\">")
	buff.WriteString(m.Status.IconWithText())
	desc := m.Description
	if desc == "" {
		desc = m.Operation
	}
	buff.WriteString(fmt.Sprintf("<p class=\"mutation-description\">%s</p>", html.EscapeString(desc)))
	buff.WriteString("<div class=\"spacer\"></div><div class=\"mutation-options\">")
	buff.WriteString("<button class=\"review-btn option-btn\"><img class=\"icon\" src=\"/resources/icons/pen-solid.svg\" alt=\"pen icon\" />Review</button>")
	buff.WriteString(fmt.Sprintf("<a title=\"view mutation %s\" href=\"/%s/mutant/%s?m=%s#%s\">%.7s</a>", m.ID, r.framework, r.file, m.ID, m.ID, m.ID))
	buff.WriteString("</div></div></td></tr>")
}

func (r *codeRenderer) renderAllMutationData(buff *bytes.Buffer, m *mutations.Mutation) {
	buff.WriteString("<tr><td colspan=\"100%\"><div class=\"all-data-wrapper\">")

	// Marv Mutant ID
	buff.WriteString(fmt.Sprintf("<p><span class=\"data-type\">Mutant ID (Marv):</span> %s</p>", m.ID))

	// Framework Mutant ID
	buff.WriteString(fmt.Sprintf("<p><span class=\"data-type\">Mutant ID (%s):</span> ", r.framework))
	if m.FrameworkMutantID == "" {
		buff.WriteString(fmt.Sprintf("<span class=\"orange\">Framework <strong>%s</strong> does not create mutant ids</span>", r.framework))
	} else {
		buff.WriteString(m.FrameworkMutantID)
	}
	buff.WriteString("</p>")

	buff.WriteString(fmt.Sprintf("<p><span class=\"data-type\">Description:</span> %s</p>", html.EscapeString(m.Description)))

	// Mutation Operator
	buff.WriteString(fmt.Sprintf("<p><span class=\"data-type\">Mutation Operator:</span> %s</p>", html.EscapeString(m.Operation)))
}

func (r *codeRenderer) renderReviewField(buff *bytes.Buffer, m *mutations.Mutation) error {
	rev, err := r.db.GetReviewByMutationID(m.ID)
	switch true {
	case errors.Is(err, gorm.ErrRecordNotFound):
		rev = &review.Review{}
	case err != nil:
		return err
	}

	buff.WriteString("<tr class=\"review")
	if rev.Review == "" {
		buff.WriteString(" hidden")
	}
	buff.WriteString("\">")

	buff.WriteString(fmt.Sprintf("<td colspan=\"100%%\">"+
		"<div class=\"review-wrapper\">"+
		"<div class=\"review-header\">"+
		"<label for=\"review-%s\" class=\"generic-label\">Add Review</label>"+
		"<div class=\"loader-wrapper saved\">"+
		"<img class=\"saved-icon\" src=\"/resources/icons/circle-check-solid.svg\" alt=\"saved icon\" />"+
		"<div class=\"loader\"></div>"+
		"<p class=\"loader-status\">Saved</p>"+
		"</div>"+ // closes loader-wrapper
		"</div>"+ // closes review-header
		"<textarea id=\"review-%s\" class=\"generic-textarea\" type=\"text\" placeholder=\"Enter review...\">%s</textarea>"+
		"</div>"+
		"</td></tr>", m.ID, m.ID, rev.Review))
	return nil
}
