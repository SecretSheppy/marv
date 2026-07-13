package html

import (
	"bytes"

	"github.com/SecretSheppy/marv/web"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

type resourcesRenderer struct {
	shared *shared
}

func newResourcesRenderer(shared *shared) *resourcesRenderer {
	return &resourcesRenderer{shared: shared}
}

func (r *resourcesRenderer) Render(buff *bytes.Buffer) error {
	return r.minify(buff)
}

func (r *resourcesRenderer) minify(buff *bytes.Buffer) error {
	mini := minify.New()
	mini.AddFunc("text/css", css.Minify)
	mini.AddFunc("text/javascript", js.Minify)
	if err := r.minifyStyles(buff, mini); err != nil {
		return err
	}
	buff.WriteString("<style>")
	buff.WriteString(r.shared.document.Theme.CSS())
	buff.WriteString("</style>")
	return r.minifyScripts(buff, mini)
}

func (r *resourcesRenderer) minifyStyles(buff *bytes.Buffer, mini *minify.M) error {
	buff.WriteString("<style>")
	for _, style := range r.shared.document.Stylesheets {
		content, err := web.StylesFS.ReadFile(style)
		if err != nil {
			return err
		}
		minified, err := mini.Bytes("text/css", content)
		if err != nil {
			return err
		}
		buff.Write(minified)
	}
	buff.WriteString("</style>")
	return nil
}

func (r *resourcesRenderer) minifyScripts(buff *bytes.Buffer, mini *minify.M) error {
	for _, script := range r.shared.document.Scripts {
		content, err := web.ScriptsFS.ReadFile(script)
		if err != nil {
			return err
		}
		minified, err := mini.Bytes("text/javascript", content)
		if err != nil {
			return err
		}
		buff.WriteString("<script type=\"text/javascript\">")
		buff.Write(minified)
		buff.WriteString("</script>")
	}
	return nil
}
