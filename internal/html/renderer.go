package html

import (
	"bytes"
	"errors"
	"fmt"
	"path"

	"github.com/SecretSheppy/marv/fwlib"
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

type Document struct {
	Theme                *themes.Theme
	Favicon              string
	Stylesheets, Scripts []string
}

type Renderer struct {
	cache      cache
	document   *Document
	frameworks []fwlib.Framework
	db         *review.Repository
}

func NewRenderer(document *Document, frameworks []fwlib.Framework, db *review.Repository) *Renderer {
	return &Renderer{
		cache:      make(cache),
		document:   document,
		frameworks: frameworks,
		db:         db,
	}
}

func (r *Renderer) getResources() (string, error) {
	meta := r.cache.get("resources")
	if meta == "" {
		var temp bytes.Buffer
		t := &resourcesRenderer{r.document.Stylesheets, r.document.Scripts, r.document.Theme}
		if err := t.Render(&temp); err != nil {
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
		t := &treeRenderer{r.frameworks, r.document.Theme}
		t.Render(&temp)
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
	m := &metaRenderer{title, r.document.Favicon}
	m.Render(buff)
	buff.WriteString("</head><body>")
	return nil
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
	writeFilters(&buff, r.document.Theme)
	buff.WriteString("</div>")

	buff.WriteString("<div class=\"content-wrapper\"><div class=\"content-header\">")
	buff.WriteString(fmt.Sprintf("<img class=\"content-icon\" src=\"%s\" alt=\"chart icon\" />",
		r.document.Theme.Icon("chart-simple-solid.svg")))
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

	for _, framework := range r.frameworks {
		meta := framework.Meta()
		for f, _ := range framework.Mutations() {
			stats := framework.Mutations().StatisticsFrom(f)
			buff.WriteString("<tr>")
			buff.WriteString(fmt.Sprintf("<td><a href=\"/%s/mutants/%s\">"+
				"<img class=\"icon\" src=\"%s\" alt=\"%s language icon\"/>%s</a></td>",
				meta.Name, f, meta.Language.Icon(), meta.Language.Name(), f))
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

func (r *Renderer) renderCode(framework fwlib.Framework, filePath string, conflicts mutations.Conflicts, config *codeRendererConfig) ([]byte, string, error) {
	lines, err := framework.ReadLines(filePath)
	if err != nil {
		return nil, "", err
	}
	meta := framework.Meta()
	c, err := newCodeRenderer(meta.Language.Ext(), meta.Name, filePath, lines, conflicts, config, r.db, r.document.Theme)
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
	lang := meta.Language

	var buff bytes.Buffer
	crConfig := &codeRendererConfig{RenderAllMutantData: config.Features.AdvancedDetail}
	render, codeStyle, err := r.renderCode(config.Framework, config.FilePath, config.conflicts(), crConfig)
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
	writeFilters(&buff, r.document.Theme)
	buff.WriteString("</div>") // closes sidebar-wrapper

	buff.WriteString("<div class=\"content-wrapper\"><div class=\"content-header\">")
	writeFrameworkName(&buff, config.Framework)
	buff.WriteString(fmt.Sprintf("<img class=\"content-icon\" src=\"%s\" alt=\"%s language icon\" />"+
		"<h3 class=\"content-title\">%s</h3></div>", lang.Icon(), lang.Name(), path.Base(config.FilePath)))
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
