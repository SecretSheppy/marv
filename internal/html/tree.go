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

type NodeType int

const (
	Directory NodeType = iota
	File
)

// PathNode represents a node in the file tree. A node can be either a Directory or a File.
type PathNode struct {
	Type     NodeType
	Name     string
	children []*PathNode
}

func (p *PathNode) AddFile(filePath string) {
	split := strings.Split(filePath, string(os.PathSeparator))
	if len(split) > 1 {
		node := p.ChildNode(split[0])
		if node == nil {
			node = &PathNode{Type: Directory, Name: split[0]}
			p.children = append(p.children, node)
		}
		node.AddFile(path.Join(split[1:]...))
	} else {
		node := &PathNode{Type: File, Name: split[0]}
		p.children = append(p.children, node)
	}
}

func (p *PathNode) SortChildren() {
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

func (p *PathNode) Children() []*PathNode {
	return p.children
}

func (p *PathNode) ChildNode(name string) *PathNode {
	for _, child := range p.children {
		if child.Name == name {
			return child
		}
	}
	return nil
}

func (p *PathNode) FirstChild() *PathNode {
	if len(p.children) == 0 {
		return nil
	}
	return p.children[0]
}

func (p *PathNode) MergeOnlyChildren() {
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

func (p *PathNode) Render(buff *bytes.Buffer, fw fwlib.Framework) {
	p.render(buff, fw, fw.Meta().Language, 1, fmt.Sprintf("/%s/mutants", fw.Meta().Name))
}

func (p *PathNode) render(buff *bytes.Buffer, fw fwlib.Framework, lang *languages.Language, level int, accPath string) {
	currentPath := path.Join(accPath, p.Name)
	switch p.Type {
	case Directory:
		buff.WriteString("<div class=\"directory-wrapper collapsed\">")
		p.renderDirectoryNode(buff, level)
		buff.WriteString("<div class=\"directory-contents\">")
		for _, child := range p.children {
			child.render(buff, fw, lang, level+1, currentPath)
		}
		buff.WriteString("</div></div>")
	case File:
		p.renderFileNode(buff, level, currentPath, fw.Meta().Language)
	}
}

func (p *PathNode) renderDirectoryNode(buff *bytes.Buffer, level int) {
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
		"</div>"+ // closes icon-name-wrapper
		"</div>", level, p.Name))
}

func (p *PathNode) renderFileNode(buff *bytes.Buffer, level int, href string, lang *languages.Language) {
	buff.WriteString(fmt.Sprintf("<a class=\"node file\" style=\"--level: %d;\" href=\"%s\">"+
		"<div class=\"spacer\"></div>"+
		"<div class=\"icon-name-wrapper\">"+
		"<img class=\"icon\" src=\"%s\" alt=\"%s language icon\" />"+
		"<p class=\"name\">%s</p>"+
		"</div>"+ // closes icon-name-wrapper
		"</a>", level, href, lang.Icon(), lang.Name(), p.Name))
}

type TreeRenderer struct {
	fws []fwlib.Framework
}

func NewTreeRenderer(fws []fwlib.Framework) *TreeRenderer {
	return &TreeRenderer{fws: fws}
}

func (t *TreeRenderer) Render(buff *bytes.Buffer) {
	buff.WriteString("<div class=\"tree\">")
	t.renderHeader(buff)
	buff.WriteString("<div class=\"tree-body\">")
	for _, fw := range t.fws {
		buff.WriteString("<div class=\"framework\">")
		t.renderFrameworkHeader(buff, fw)
		root := &PathNode{}
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

func (t *TreeRenderer) renderHeader(buff *bytes.Buffer) {
	buff.WriteString("<div class=\"tree-header\">" +
		"<a href=\"/start\"><img class=\"header-logo\" src=\"/resources/branding/marv_logo.png\" alt=\"marv logo\" /></a>" +
		"<div class=\"buttons-wrapper\">" +
		"<button class=\"header-button\" title=\"Locate Selected File\"><img class=\"icon\" src=\"/resources/icons/crosshair.png\" alt=\"crosshair icon\" /></button>" +
		"<button class=\"header-button\" title=\"Expand Selected\"><img class=\"icon\" src=\"/resources/icons/up_down.png\" alt=\"up arrow above down arrow icon\" /></button>" +
		"<button class=\"header-button\" title=\"Collapse All\"><img class=\"icon\" src=\"/resources/icons/down_up.png\" alt=\"down arrow above up arrow icon\" /></button>" +
		"</div>" +
		"</div>")
}

func (t *TreeRenderer) renderFrameworkHeader(buff *bytes.Buffer, fw fwlib.Framework) {
	// TODO: refactor this into the mutations package
	var total, covered, killed float64
	for _, conflicts := range fw.Mutations() {
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
	coverage := covered / total * 100
	score := killed / total * 100
	ofCovered := killed / covered * 100
	buff.WriteString(fmt.Sprintf("<div class=\"framework-header\">"+
		"<div class=\"framework-name\">%s</div>"+
		"<p class=\"stats\">coverage: <span class=\"%s\">%.2f%%</span>,</p>"+
		"<p class=\"stats\">score: <span class=\"%s\">%.2f%%</span>,</p>"+
		"<p class=\"stats\">of covered: <span class=\"%s\">%.2f%%</span></p>"+
		"</div>", fw.Meta().Name, statGradeClass(coverage), coverage,
		statGradeClass(score), score, statGradeClass(ofCovered), ofCovered))
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
