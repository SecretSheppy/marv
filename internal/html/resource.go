package html

import (
	"bytes"
	"fmt"
	"io/fs"

	"github.com/SecretSheppy/marv/web"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

func minit() *minify.M {
	mini := minify.New()
	mini.AddFunc("text/css", css.Minify)
	mini.AddFunc("text/javascript", js.Minify)
	return mini
}

type sources struct {
	styles, scripts fs.FS
}

func defaultSources() *sources {
	return &sources{styles: web.StylesFS, scripts: web.ScriptsFS}
}

type resourcesRenderer struct {
	shared   *shared
	minifier *minify.M
	sources  *sources
}

func newResourcesRenderer(shared *shared) *resourcesRenderer {
	return &resourcesRenderer{shared: shared, minifier: minit(), sources: defaultSources()}
}

func (r *resourcesRenderer) style(buff *bytes.Buffer, stylesheet string) error {
	content, err := fs.ReadFile(r.sources.styles, stylesheet)
	if err != nil {
		return err
	}
	minified, err := r.minifier.Bytes("text/css", content)
	if err != nil {
		return err
	}
	buff.Write(minified)
	return nil
}

func (r *resourcesRenderer) styles(buff *bytes.Buffer) error {
	buff.WriteString("<style>")
	for _, stylesheet := range r.shared.document.Stylesheets {
		if err := r.style(buff, stylesheet); err != nil {
			return err
		}
	}
	buff.WriteString(r.shared.document.Theme.CSS())
	buff.WriteString("</style>")
	return nil
}

func (r *resourcesRenderer) script(buff *bytes.Buffer, script string) error {
	content, err := fs.ReadFile(r.sources.scripts, script)
	if err != nil {
		return err
	}
	minified, err := r.minifier.Bytes("text/javascript", content)
	if err != nil {
		return err
	}
	buff.WriteString(fmt.Sprintf(`<script type="text/javascript">%s</script>`, minified))
	return nil
}

func (r *resourcesRenderer) scripts(buff *bytes.Buffer) error {
	for _, script := range r.shared.document.Scripts {
		r.script(buff, script)
	}
	return nil
}

func (r *resourcesRenderer) RenderResources(buff *bytes.Buffer) error {
	return r.render(buff)
}

func (r *resourcesRenderer) render(buff *bytes.Buffer) error {
	if err := r.styles(buff); err != nil {
		return err
	}
	return r.scripts(buff)
}

func (r *resourcesRenderer) RenderScripts(buff *bytes.Buffer) error {
	return r.scripts(buff)
}
