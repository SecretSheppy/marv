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
	"github.com/SecretSheppy/marv/pkg/chroma_proxy"
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

type codeRenderer struct {
	shared             *shared
	config             *RenderConfig
	lines, highlighted []string
	lnPadding          int // automatically configured value used to align line numbers throughout the document.
	proxy              *chroma_proxy.ProxyHighlighter
}

func newCodeRenderer(shared *shared, config *RenderConfig) *codeRenderer {
	return &codeRenderer{shared: shared, config: config}
}

func (r *codeRenderer) Highlighter() *chroma_proxy.ProxyHighlighter {
	return r.proxy
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
	conflicts := r.config.conflicts()
	conflicts.Sort()
	rendered := make([]*renderedConflict, 0, len(conflicts))

	// TODO: move out of this function
	lines, err := r.config.lines()
	if err != nil {
		return nil, err
	}
	r.lines = lines

	// TODO: move out of this function
	proxy, err := chroma_proxy.NewProxyHighlighter(r.config.language().MExt(), r.shared.document.Theme.Code.ChromaTheme)
	if err != nil {
		return nil, err
	}
	r.proxy = proxy

	highlighted, err := proxy.Highlight(r.lines)
	if err != nil {
		return nil, err
	}
	r.highlighted = highlighted

	for _, conflict := range conflicts {
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

		r.renderLine(&buff, i+1, lineEqual, r.highlighted[i])
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
	r.renderDebugMessage(&buff, c.String())
	buff.WriteString(fmt.Sprintf("<tbody class=\"hidden\" data-conflict-id=\"%s\">", c.ID))
	r.renderLines(&buff, c.StartLine, c.EndLine)
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

func (r *codeRenderer) renderLines(buff *bytes.Buffer, from, to int) {
	for i, line := range r.highlighted[from : to+1] {
		r.renderLine(buff, from+i+1, lineEqual, line)
	}
}

func (r *codeRenderer) renderMutation(c *mutations.Conflict, m *mutations.Mutation) (string, error) {
	var buff bytes.Buffer
	r.renderDebugMessage(&buff, m.String())
	buff.WriteString(fmt.Sprintf("<tbody id=\"%s\" data-conflict-id=\"%s\" data-status=\"%s\" data-class=\"mutant\" class=\"mutation\">", m.ID, c.ID, m.Status.Text()))
	r.renderMutationHeader(&buff, m)
	if r.config.Features.AdvancedDetail {
		r.renderAllMutationData(&buff, m)
	}

	r.renderLines(&buff, c.StartLine, m.Start.Line-1)

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
		//lines, err := r.proxy.Highlight([]string{pre, diff, post})
		//if err != nil {
		//	return "", err
		//}
		lines := []string{pre, diff, post}
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
		//lines, err := r.proxy.Highlight([]string{pre, diff, post})
		//if err != nil {
		//	return "", err
		//}
		lines := []string{pre, diff, post}
		code := fmt.Sprintf("%s<span class=\"highlight insert\">%s</span>%s", lines[0], lines[1], lines[2])
		r.renderLine(&buff, 0, lineInserted, code)
	}

	r.renderLines(&buff, m.End.Line+1, c.EndLine)

	if err := r.renderReviewField(&buff, m); err != nil {
		return "", err
	}

	buff.WriteString("</tbody>")
	return buff.String(), nil
}

func (r *codeRenderer) renderMutationHeader(buff *bytes.Buffer, m *mutations.Mutation) {
	buff.WriteString("<tr><td colspan=\"100%\"><div class=\"mutation-header\">")
	buff.WriteString(m.Status.IconWithText())
	buff.WriteString(fmt.Sprintf("<p class=\"mutation-description\">%s</p>", html.EscapeString(m.GetDescription())))
	buff.WriteString("<div class=\"spacer\"></div><div class=\"mutation-options\">")
	buff.WriteString("<button class=\"review-btn option-btn\"><img class=\"icon\" src=\"" + getIconURL(r.shared.document.Theme, "pen-solid.svg") + "\" alt=\"pen icon\" />Review</button>")
	buff.WriteString(fmt.Sprintf("<a title=\"view mutation %s\" href=\"/%s/mutant/%s?m=%s#%s\">%.7s</a>", m.ID, r.config.Framework.Meta().Name, r.config.FilePath, m.ID, m.ID, m.ID))
	buff.WriteString("</div></div></td></tr>")
}

func (r *codeRenderer) renderAllMutationData(buff *bytes.Buffer, m *mutations.Mutation) {
	buff.WriteString("<tr><td colspan=\"100%\"><div class=\"all-data-wrapper\">")

	// Marv Mutant ID
	buff.WriteString(fmt.Sprintf("<p><span class=\"data-type\">Mutant ID (Marv):</span> %s</p>", m.ID))

	// Framework Mutant ID
	buff.WriteString(fmt.Sprintf("<p><span class=\"data-type\">Mutant ID (%s):</span> ", r.config.Framework.Meta().Name))
	if m.FrameworkMutantID == "" {
		buff.WriteString(fmt.Sprintf("<span class=\"orange\">Framework <strong>%s</strong> does not create mutant ids</span>", r.config.Framework.Meta().Name))
	} else {
		buff.WriteString(m.FrameworkMutantID)
	}
	buff.WriteString("</p>")

	buff.WriteString(fmt.Sprintf("<p><span class=\"data-type\">Description:</span> %s</p>", html.EscapeString(m.Description)))

	// Mutation Operator
	buff.WriteString(fmt.Sprintf("<p><span class=\"data-type\">Mutation Operator:</span> %s</p>", html.EscapeString(m.Operation)))
}

func (r *codeRenderer) renderReviewField(buff *bytes.Buffer, m *mutations.Mutation) error {
	rev, err := r.shared.db.GetReviewByMutationID(m.ID)
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
		"<img class=\"saved-icon\" src=\"%s\" alt=\"saved icon\" />"+
		"<div class=\"loader\"></div>"+
		"<p class=\"loader-status\">Saved</p>"+
		"</div>"+ // closes loader-wrapper
		"</div>"+ // closes review-header
		"<textarea id=\"review-%s\" class=\"generic-textarea\" type=\"text\" placeholder=\"Enter review...\">%s</textarea>"+
		"</div>"+
		"</td></tr>", m.ID, r.shared.document.Theme.Icon("circle-check-solid.svg"), m.ID, rev.Review))
	return nil
}

func (r *codeRenderer) renderDebugMessage(buff *bytes.Buffer, msg string) {
	if !r.shared.debug {
		return
	}
	buff.WriteString(fmt.Sprintf(`<tbody><tr><td colspan="100%%" style="background:orange;color:black">%s</td></tr></tbody>`, msg))
}
