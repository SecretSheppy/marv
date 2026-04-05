package html

import (
	"bytes"
	"io"

	"github.com/SecretSheppy/marv/internal/mutations"
)

// Renderer produces the entire HTML output that is displayed to the user. It does this by calling and combining the
// results of several sub-renderer objects.
type Renderer struct {
	meta *Meta
	code *CodeRenderer
}

func NewRenderer(meta *Meta, ext string, lines []string, conflicts mutations.Conflicts) (*Renderer, error) {
	code, err := NewCodeRenderer(ext, lines, conflicts)
	if err != nil {
		return nil, err
	}
	return &Renderer{meta: meta, code: code}, nil
}

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
	var buff bytes.Buffer

	r.scripts(&buff)
	if err := r.styles(&buff); err != nil {
		return nil, err
	}
	if err := r.code.Render(&buff); err != nil {
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

	code, err := r.code.SyntaxHighlighter().CSS()
	if err != nil {
		return err
	}
	w.WriteString(code)

	w.WriteString("</style>")
	return nil
}
