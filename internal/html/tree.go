package html

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mutations"
)

type nodeType int

const (
	directory nodeType = iota
	file
)

// pathNode represents a node in the file tree. A node can be either a directory or a file.
type pathNode struct {
	Name     string
	Type     nodeType
	children []*pathNode
}

func (p *pathNode) AddFile(filePath string) {
	parts := strings.Split(filePath, string(os.PathSeparator))
	if len(parts) == 1 {
		p.children = append(p.children, &pathNode{Type: file, Name: parts[0]})
		return
	}

	node := p.ChildNode(parts[0])
	if node == nil {
		node = &pathNode{Type: directory, Name: parts[0]}
		p.children = append(p.children, node)
	}
	node.AddFile(path.Join(parts[1:]...))
}

func (p *pathNode) SortChildren() {
	sort.Slice(p.children, func(i, j int) bool {
		if p.children[i].Type != p.children[j].Type {
			return p.children[i].Type < p.children[j].Type
		}
		return p.children[i].Name < p.children[j].Name
	})
	for _, child := range p.children {
		child.SortChildren()
	}
}

func (p *pathNode) Children() []*pathNode {
	return p.children
}

func (p *pathNode) FirstChild() *pathNode {
	if len(p.children) == 0 {
		return nil
	}
	return p.children[0]
}

func (p *pathNode) ChildNode(name string) *pathNode {
	for _, child := range p.children {
		if child.Name == name {
			return child
		}
	}
	return nil
}

func (p *pathNode) MergeOnlyChildren() {
	for _, child := range p.children {
		child.MergeOnlyChildren()
	}

	if len(p.children) == 1 {
		firstChild := p.children[0]
		p.Name = fmt.Sprintf("%s/%s", p.Name, firstChild.Name)
		p.Type = firstChild.Type
		p.children = firstChild.children
	}
}

func (p *pathNode) Render(buff *bytes.Buffer, fw fwlib.Framework) {
	p.render(buff, fw, fw.Meta().Language, 1, "")
}

func (p *pathNode) render(buff *bytes.Buffer, fw fwlib.Framework, lang *languages.Language, level int, accPath string) {
	currentPath := path.Join(accPath, p.Name)
	switch p.Type {
	case directory:
		buff.WriteString("<div class=\"directory-wrapper collapsed\">")
		p.renderDirectoryNode(buff, level, currentPath, fw)
		buff.WriteString("<div class=\"directory-contents\">")
		for _, child := range p.children {
			child.render(buff, fw, lang, level+1, currentPath)
		}
		buff.WriteString("</div></div>")
	case file:
		p.renderFileNode(buff, level, currentPath, fw)
	}
}

func (p *pathNode) renderDirectoryNode(buff *bytes.Buffer, level int, currentPath string, fw fwlib.Framework) {
	buff.WriteString(fmt.Sprintf("<div class=\"node directory\" style=\"--level: %d;\">"+
		"<div class=\"spacer\">"+
		"<div class=\"collapse-toggle\">"+
		"<img class=\"icon icon-expanded\" src=\"/resources/icons/arrow_down.png\" alt=\"down arrow\" />"+
		"<img class=\"icon icon-collapsed\" src=\"/resources/icons/arrow_right.png\" alt=\"right arrow\" />"+
		"</div>"+ // closes collapse-toggle
		"</div>"+ // closes spacer
		"<div class=\"icon-name-wrapper\">"+
		"<img class=\"icon\" src=\"/resources/icons/folder-solid.svg\" alt=\"folder icon\" />"+
		"<p class=\"name\">%s</p>"+
		"</div>", level, p.Name))
	writeWrappedStats(buff, currentPath, fw)
	buff.WriteString("</div>")
}

func (p *pathNode) renderFileNode(buff *bytes.Buffer, level int, href string, fw fwlib.Framework) {
	lang := fw.Meta().Language
	prefix := fmt.Sprintf("/%s/mutants/", fw.Meta().Name)
	buff.WriteString(fmt.Sprintf("<a class=\"node file\" style=\"--level: %d;\" href=\"%s%s\">"+
		"<div class=\"spacer\"></div>"+
		"<div class=\"icon-name-wrapper\">"+
		"<img class=\"icon\" src=\"%s\" alt=\"%s language icon\" />"+
		"<p class=\"name\">%s</p>"+
		"</div>", level, prefix, href, lang.Icon(), lang.Name(), p.Name))
	writeWrappedStats(buff, href, fw)
	buff.WriteString("</a>")
}

type treeRenderer struct {
	fws []fwlib.Framework
}

func (t *treeRenderer) Render(buff *bytes.Buffer) {
	buff.WriteString("<div id=\"fw-tree\" class=\"tree\">")
	t.renderHeader(buff)
	buff.WriteString("<div id=\"tree-body\" class=\"tree-body\">")
	for _, fw := range t.fws {
		buff.WriteString("<div class=\"framework\">")
		t.renderFrameworkHeader(buff, fw)
		root := &pathNode{}
		for k, _ := range fw.Mutations() {
			root.AddFile(k)
		}
		root.FirstChild().MergeOnlyChildren()
		root.SortChildren()
		root.FirstChild().Render(buff, fw)
		buff.WriteString("</div>")
	}
	buff.WriteString("</div></div>")
}

func (t *treeRenderer) renderHeader(buff *bytes.Buffer) {
	buff.WriteString("<div class=\"tree-header\">" +
		"<a href=\"/start\"><img class=\"header-logo\" src=\"/resources/branding/marv_logo.png\" alt=\"marv logo\" /></a>" +
		"<div class=\"buttons-wrapper\">" +
		"<button id=\"tree-crosshair-btn\" class=\"header-button\" title=\"Locate Selected File\"><img class=\"icon\" src=\"/resources/icons/crosshair.png\" alt=\"crosshair icon\" /></button>" +
		"<button id=\"tree-expand-all-btn\" class=\"header-button\" title=\"Expand All\"><img class=\"icon\" src=\"/resources/icons/up_down.png\" alt=\"up arrow above down arrow icon\" /></button>" +
		"<button id=\"tree-collapse-all-btn\" class=\"header-button\" title=\"Collapse All\"><img class=\"icon\" src=\"/resources/icons/down_up.png\" alt=\"down arrow above up arrow icon\" /></button>" +
		"</div>" +
		"</div>")
}

func (t *treeRenderer) renderFrameworkHeader(buff *bytes.Buffer, fw fwlib.Framework) {
	buff.WriteString("<div class=\"framework-header\">")
	writeFrameworkName(buff, fw)
	writeStats(buff, "", fw)
	buff.WriteString("</div>")
}

func writeFrameworkName(buff *bytes.Buffer, fw fwlib.Framework) {
	meta := fw.Meta()
	buff.WriteString(fmt.Sprintf("<div class=\"framework-name\" "+
		"title=\"Mutants created by the %s mutation testing framework for %s\">%s</div>",
		meta.Name, meta.Language.Name(), meta.Name))
}

func writeWrappedStats(buff *bytes.Buffer, startPath string, fw fwlib.Framework) {
	buff.WriteString("<div class=\"stats-wrapper\">")
	writeStats(buff, startPath, fw)
	buff.WriteString("</div>")
}

func writeStats(buff *bytes.Buffer, startPath string, fw fwlib.Framework) {
	var total, covered, killed float64
	for mutPath, conflicts := range fw.Mutations() {
		if strings.HasPrefix(mutPath, startPath) {
			for _, conflict := range conflicts {
				for _, mutation := range conflict.Mutations {
					total++
					if mutation.Status != mutations.NoCoverage {
						covered++
					}
					if mutation.Status == mutations.Killed {
						killed++
					}
				}
			}
		}
	}
	coverage := covered / total * 100
	score := killed / total * 100
	ofCovered := killed / covered * 100
	buff.WriteString(fmt.Sprintf("<p class=\"stats\">coverage: <span class=\"%s\">%.2f%%</span>,</p>",
		statGradeClass(coverage), coverage))
	buff.WriteString(fmt.Sprintf("<p class=\"stats\">score: <span class=\"%s\">%.2f%%</span>,</p>",
		statGradeClass(score), score))
	buff.WriteString(fmt.Sprintf("<p class=\"stats\">of covered: <span class=\"%s\">%.2f%%</span></p>",
		statGradeClass(ofCovered), ofCovered))
}

func statGradeClass(stat float64) string {
	if stat >= 80 {
		return "green"
	}
	if stat >= 60 {
		return "orange"
	}
	return "red"
}
