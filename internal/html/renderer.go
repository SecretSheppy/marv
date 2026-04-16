package html

import (
	"bytes"
	"errors"
	"fmt"
	"path"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/SecretSheppy/marv/pkg/fio"
	"github.com/google/uuid"
)

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

type Config struct {
	Favicon         string
	Styles, Scripts []string
}

type Renderer struct {
	cache      cache
	config     *Config
	frameworks []fwlib.Framework
}

func NewRenderer(config *Config, frameworks []fwlib.Framework) *Renderer {
	return &Renderer{
		cache:      make(cache),
		config:     config,
		frameworks: frameworks,
	}
}

func (r *Renderer) getResources() (string, error) {
	meta := r.cache.get("resources")
	if meta == "" {
		var temp bytes.Buffer
		t := &resourcesRenderer{r.config.Styles, r.config.Scripts}
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
		t := &treeRenderer{r.frameworks}
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
	m := &metaRenderer{title, r.config.Favicon}
	m.Render(buff)
	buff.WriteString("</head><body>")
	return nil
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

func (r *Renderer) renderCode(framework fwlib.Framework, filePath string, conflicts mutations.Conflicts) ([]byte, string, error) {
	absolutePath := path.Join(framework.Yaml().SourceCodeDir(), filePath)
	lines, err := fio.ReadLines(absolutePath)
	if err != nil {
		return nil, "", err
	}
	meta := framework.Meta()
	c, err := newCodeRenderer(meta.Language.Ext(), meta.Name, filePath, lines, conflicts)
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

func (r *Renderer) renderMutants(framework fwlib.Framework, conflicts mutations.Conflicts, filePath, title string) ([]byte, error) {
	var buff bytes.Buffer
	render, codeStyle, err := r.renderCode(framework, filePath, conflicts)
	if err != nil {
		return nil, err
	}
	if err := r.renderHead(&buff, title, "<style>"+codeStyle+"</style>"); err != nil {
		return nil, err
	}
	buff.WriteString(r.getTree())
	buff.Write(render)
	buff.WriteString("</body></html>")
	return buff.Bytes(), nil
}

func (r *Renderer) RenderMutant(framework fwlib.Framework, filePath string, mutantID uuid.UUID) ([]byte, error) {
	title := fmt.Sprintf("[%s] %s -> mutant[%s]", framework, filePath, mutantID)

	conflict, mutant := framework.Mutations()[filePath].GetMutant(mutantID)
	if mutant == nil {
		return nil, errors.New("mutant not found with id " + mutantID.String())
	}
	conflicts := mutations.Conflicts{
		&mutations.Conflict{
			ID:        conflict.ID,
			StartLine: conflict.StartLine,
			EndLine:   conflict.EndLine,
			Mutations: []*mutations.Mutation{mutant},
		},
	}

	return r.renderMutants(framework, conflicts, filePath, title)
}

func (r *Renderer) RenderMutants(framework fwlib.Framework, filePath string) ([]byte, error) {
	title := fmt.Sprintf("[%s] %s", framework, filePath)
	return r.renderMutants(framework, framework.Mutations()[filePath], filePath, title)
}
