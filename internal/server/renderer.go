package server

import (
	"bytes"
	"io"

	"github.com/SecretSheppy/marv/pkg/highlighter"
	"github.com/alecthomas/chroma/v2/styles"
)

type Renderer struct {
	ext       string
	lines     []string
	highlight *highlighter.Highlighter
	meta      *Meta
}

func NewRenderer(ext string, lines []string, meta *Meta) *Renderer {
	return &Renderer{ext: ext, lines: lines, meta: meta}
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

func (r *Renderer) code(w *bytes.Buffer) error {
	lines, err := r.highlight.HighlightLines(0, len(r.lines)-1)
	if err != nil {
		return err
	}

	w.WriteString("<table>")

	for _, line := range lines {
		w.WriteString("<tr><td>")
		w.WriteString(line)
		w.WriteString("</td></tr>")
	}

	w.WriteString("</table>")
	return nil
}
