package html

import (
	"bytes"
	"errors"
	"fmt"
	"path"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/internal/review"
	"github.com/SecretSheppy/marv/internal/themes"
	"github.com/google/uuid"
)

func getIconURL(theme *themes.Theme, name string) string {
	return fmt.Sprintf("/icon/%s/%s", theme.IconColor(), name)
}

type cache map[string]string

func (c cache) set(name, content string) {
	c["data:"+name] = content
}

func (c cache) get(name string) string {
	return c["data:"+name]
}

func (c cache) setFile(file, content string) {
	c["file:"+file] = content
}

func (c cache) getFile(file string) string {
	return c["file:"+file]
}

type RenderFeatures struct {
	Filtering, AdvancedDetail bool
}

type RenderConfig struct {
	Framework fwlib.Framework
	Conflicts mutations.Conflicts
	FilePath  string
	Features  *RenderFeatures
}

func (r *RenderConfig) title() string {
	return fmt.Sprintf("[%s] %s", r.Framework.Meta().Name, r.FilePath)
}

func (r *RenderConfig) conflicts() mutations.Conflicts {
	if r.Conflicts != nil {
		return r.Conflicts
	}
	return r.Framework.Mutations()[r.FilePath]
}

func (r *RenderConfig) lines() ([]string, error) {
	return r.Framework.ReadLines(r.FilePath)
}

func (r *RenderConfig) language() *languages.Language {
	return languages.GetLanguageFromFile(r.FilePath)
}

type Document struct {
	Theme                *themes.Theme
	Favicon              string
	Stylesheets, Scripts []string
}

type shared struct {
	db         *review.Repository
	document   *Document
	frameworks []fwlib.Framework
}

type Renderer struct {
	cache  cache
	shared *shared
}

func NewRenderer(document *Document, db *review.Repository, frameworks []fwlib.Framework) *Renderer {
	return &Renderer{
		cache:  make(cache),
		shared: &shared{db: db, document: document, frameworks: frameworks},
	}
}

func (r *Renderer) getResources() (string, error) {
	meta := r.cache.get("resources")
	if meta == "" {
		var temp bytes.Buffer
		if err := newResourcesRenderer(r.shared).Render(&temp); err != nil {
			return "", err
		}
		meta = temp.String()
		r.cache.set("resources", meta)
	}
	return meta, nil
}

func (r *Renderer) getTree() string {
	tree := r.cache.get("tree")
	if tree == "" {
		var temp bytes.Buffer
		newTreeRenderer(r.shared).Render(&temp)
		tree = temp.String()
		r.cache.set("tree", tree)
	}
	return tree
}

func (r *Renderer) renderHead(buff *bytes.Buffer, title string, elements ...string) error {
	buff.WriteString("<!DOCTYPE html><html lang=\"en-GB\"><head><meta charset=\"UTF-8\">")
	resources, err := r.getResources()
	if err != nil {
		return err
	}
	buff.WriteString(resources)
	for _, element := range elements {
		buff.WriteString(element)
	}
	m := &metaRenderer{title, r.shared.document.Favicon}
	m.Render(buff)
	buff.WriteString("</head><body>")
	return nil
}

func (r *Renderer) renderFilters(buff *bytes.Buffer) {
	buff.WriteString("<div id=\"filters\" class=\"filters-component collapsed\">" +
		"<div id=\"filters-toggle\" class=\"filters-bar\">" +
		"<img class=\"icon\" src=\"" + r.shared.document.Theme.Icon("sliders-solid.svg") + "\" alt=\"filters icon\" />" +
		"<h4 class=\"bar-title\">Status Filters</h4>" +
		"<div class=\"right-content\">" +
		"<img class=\"icon arrow-up\" src=\"" + r.shared.document.Theme.Icon("arrow-up.svg") + "\" alt=\"arrow up icon\" />" +
		"<img class=\"icon arrow-down\" src=\"" + r.shared.document.Theme.Icon("arrow-down.svg") + "\" alt=\"arrow down icon\" />" +
		"</div>" + // closes right-content
		"</div>" + // closes filters-bar
		"<div class=\"content-wrapper\">" +
		"<p class=\"section-description\">" +
		"Using the filters will hide all mutants with statuses that are not enabled. " +
		"This setting syncs across all open tabs." +
		"</p>" +
		"<div class=\"filters-wrapper\">")
	for _, status := range mutations.Statuses {
		buff.WriteString(fmt.Sprintf("<label for=\"show-%s\" class=\"filter\">"+
			"<input id=\"show-%s\" type=\"checkbox\" name=\"%s\" checked /> %s"+
			"</label>", status.Text(), status.Text(), status.Text(), status.Text()))
	}
	buff.WriteString("</div>" + // closes filters-wrapper
		"</div>" + // closes content-wrapper
		"</div>")
}

func (r *Renderer) RenderStart() ([]byte, error) {
	title := "Results Overview"
	var buff bytes.Buffer
	if err := r.renderHead(&buff, title, "<meta name=\"filtering-enabled\" content=\"%v\">"); err != nil {
		return nil, err
	}

	buff.WriteString("<div class=\"layout\">")
	buff.WriteString("<div class=\"sidebar-wrapper\">")
	buff.WriteString(r.getTree())
	r.renderFilters(&buff)
	buff.WriteString("</div>")

	buff.WriteString("<div class=\"content-wrapper\"><div class=\"content-header\">")
	buff.WriteString(fmt.Sprintf("<img class=\"content-icon\" src=\"%s\" alt=\"chart icon\" />",
		r.shared.document.Theme.Icon("chart-simple-solid.svg")))
	buff.WriteString(fmt.Sprintf("<h3 class=\"content-title\">%s</h3></div>", title))

	buff.WriteString("<div class=\"overflow-wrapper\"><table class=\"generic-table\">")
	buff.WriteString("<tr>" +
		"<th>File</th>" +
		"<th>Coverage</th>" +
		"<th>Score</th>" +
		"<th>Score of Covered</th>")
	for _, status := range mutations.Statuses {
		buff.WriteString(fmt.Sprintf("<th>%s</th>", status.Text()))
	}
	buff.WriteString("</tr>")

	for _, framework := range r.shared.frameworks {
		meta := framework.Meta()
		for f, _ := range framework.Mutations() {
			lang := languages.GetLanguageFromFile(f)
			stats := framework.Mutations().StatisticsFrom(f)
			buff.WriteString("<tr>")
			buff.WriteString(fmt.Sprintf("<td><a href=\"/%s/mutants/%s\">"+
				"<img class=\"icon\" src=\"%s\" alt=\"%s language icon\"/>%s</a></td>",
				meta.Name, f, lang.Icon(), lang.Name(), f))
			buff.WriteString(fmt.Sprintf("<td>%s</td>", formatColouredStat(stats.Coverage(), 2)))
			buff.WriteString(fmt.Sprintf("<td>%s</td>", formatColouredStat(stats.Score(), 2)))
			buff.WriteString(fmt.Sprintf("<td>%s</td>", formatColouredStat(stats.ScoreOfCovered(), 2)))
			for _, status := range mutations.Statuses {
				buff.WriteString(fmt.Sprintf("<td>%.0f</td>", stats.StatusCounts[status]))
			}
			buff.WriteString("</tr>")
		}
	}

	buff.WriteString("</table></div></div></body></html>")

	return buff.Bytes(), nil
}

func (r *Renderer) RenderTree() ([]byte, error) {
	title := "Frameworks Tree"
	var buff bytes.Buffer
	if err := r.renderHead(&buff, title); err != nil {
		return nil, err
	}
	buff.WriteString(r.getTree())
	buff.WriteString("</body></html>")
	return buff.Bytes(), nil
}

func (r *Renderer) renderCode(config *RenderConfig) ([]byte, string, error) {
	c, err := newCodeRenderer(r.shared, config)
	if err != nil {
		return nil, "", err
	}
	var temp bytes.Buffer
	if err := c.Render(&temp); err != nil {
		return nil, "", err
	}
	css, err := c.SyntaxHighlighter().CSS()
	if err != nil {
		return nil, "", err
	}
	return temp.Bytes(), css, nil
}

func (r *Renderer) renderMutants(config *RenderConfig) ([]byte, error) {
	meta := config.Framework.Meta()

	var buff bytes.Buffer
	render, codeStyle, err := r.renderCode(config)
	if err != nil {
		return nil, err
	}
	err = r.renderHead(&buff, config.title(),
		"<style>"+codeStyle+"</style>",
		fmt.Sprintf("<meta name=\"filtering-enabled\" content=\"%v\">", config.Features.Filtering),
		fmt.Sprintf("<meta name=\"current-file\" content=\"/%s/mutants/%s\">", meta.Name, config.FilePath),
		fmt.Sprintf("<meta name=\"current-framework\" content=\"%s\">", meta.Name))
	if err != nil {
		return nil, err
	}

	buff.WriteString("<div class=\"layout\">")
	buff.WriteString("<div class=\"sidebar-wrapper\">")
	buff.WriteString(r.getTree())
	r.renderFilters(&buff)
	buff.WriteString("</div>") // closes sidebar-wrapper

	buff.WriteString("<div class=\"content-wrapper\"><div class=\"content-header\">")
	writeFrameworkName(&buff, config.Framework)
	buff.WriteString(fmt.Sprintf("<img class=\"content-icon\" src=\"%s\" alt=\"%s language icon\" />"+
		"<h3 class=\"content-title\">%s</h3></div>",
		config.language().Icon(), config.language().Name(), path.Base(config.FilePath)))
	buff.WriteString("<div class=\"code-wrapper\">")
	buff.Write(render)
	buff.WriteString("</div><div class=\"content-gutter\">")
	writeStats(&buff, config.FilePath, config.Framework, &statsConfig{
		Count:     true,
		Coverage:  true,
		Score:     true,
		OfCovered: true,
		Crashed:   true,
		Timeout:   true,
	})
	buff.WriteString("</div></div></div></body></html>")
	return buff.Bytes(), nil
}

func (r *Renderer) RenderMutant(config *RenderConfig, mutantID uuid.UUID) ([]byte, error) {
	conflict, mutant := config.conflicts().GetMutant(mutantID)
	if mutant == nil {
		return nil, errors.New("mutant not found with id " + mutantID.String())
	}
	config.Conflicts = mutations.Conflicts{
		&mutations.Conflict{
			ID:        conflict.ID,
			StartLine: conflict.StartLine,
			EndLine:   conflict.EndLine,
			Mutations: []*mutations.Mutation{mutant},
		},
	}

	return r.renderMutants(config)
}

func (r *Renderer) RenderMutants(config *RenderConfig) ([]byte, error) {
	return r.renderMutants(config)
}
