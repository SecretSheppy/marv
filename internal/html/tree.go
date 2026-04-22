package html

import (
	"bytes"
	"fmt"
	"math"
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
	parts := strings.Split(filePath, "/")
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
		if p.Name == "" {
			p.Name = firstChild.Name
		} else {
			p.Name = fmt.Sprintf("%s/%s", p.Name, firstChild.Name)
		}
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

var (
	nodeStats     = &statsConfig{Coverage: true, Score: true, OfCovered: true}
	fwHeaderStats = &statsConfig{Count: true, Coverage: true, Score: true, OfCovered: true}
)

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
	writeWrappedStats(buff, currentPath, fw, nodeStats)
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
	writeWrappedStats(buff, href, fw, nodeStats)
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
		root.MergeOnlyChildren()
		root.SortChildren()
		root.Render(buff, fw)
		buff.WriteString("</div>")
	}
	buff.WriteString("</div></div>")
}

func (t *treeRenderer) renderHeader(buff *bytes.Buffer) {
	buff.WriteString("<div class=\"tree-header\">" +
		"<a href=\"/\"><img class=\"header-logo\" src=\"/resources/branding/marv_logo.png\" alt=\"marv logo\" /></a>" +
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
	writeStats(buff, "", fw, fwHeaderStats)
	buff.WriteString("</div>")
}

func writeFrameworkName(buff *bytes.Buffer, fw fwlib.Framework) {
	meta := fw.Meta()
	buff.WriteString(fmt.Sprintf("<div class=\"framework-name\" "+
		"title=\"Mutants created by the %s mutation testing framework for %s\">%s</div>",
		meta.Name, meta.Language.Name(), meta.Name))
}

type statsConfig struct {
	Count, Coverage, Score, OfCovered, Crashed, Timeout bool
}

func writeWrappedStats(buff *bytes.Buffer, startPath string, fw fwlib.Framework, config *statsConfig) {
	buff.WriteString("<div class=\"stats-wrapper\">")
	writeStats(buff, startPath, fw, config)
	buff.WriteString("</div>")
}

func writeStats(buff *bytes.Buffer, startPath string, fw fwlib.Framework, config *statsConfig) {
	stats := fw.Mutations().StatisticsFrom(startPath)
	if config.Count { // Total number of mutants
		buff.WriteString(fmt.Sprintf("<p class=\"stats\">count: %.0f,</p>", stats.Count))
	}
	if config.Coverage { // Percentage of mutants that are covered by the programs test suite
		buff.WriteString(fmt.Sprintf("<p class=\"stats\">coverage: %s,</p>",
			formatColouredStat(stats.Coverage(), 2)))
	}
	if config.Score { // Percentage of mutants that were killed
		buff.WriteString(fmt.Sprintf("<p class=\"stats\">score: %s,</p>",
			formatColouredStat(stats.Score(), 2)))
	}
	if config.OfCovered { // Percentage of covered mutants that were killed
		buff.WriteString(fmt.Sprintf("<p class=\"stats\">of covered: %s,</p>",
			formatColouredStat(stats.ScoreOfCovered(), 2)))
	}
	if config.Crashed { // Total number of mutants that crashed during execution
		buff.WriteString(fmt.Sprintf("<p class=\"stats\">crashed: %.0f,</p>",
			stats.StatusCounts[mutations.Crashed]))
	}
	if config.Timeout { // Total number of mutants that timed out during execution
		buff.WriteString(fmt.Sprintf("<p class=\"stats\">timeout: %.0f</p>",
			stats.StatusCounts[mutations.Timeout]))
	}
}

func formatColouredStat(stat float64, decimalPlaces int) string {
	if math.IsNaN(stat) {
		return "<span class=\"gray\">—</span>"
	}
	var class string
	switch true {
	case stat >= 80:
		class = "green"
	case stat >= 60:
		class = "orange"
	default:
		class = "red"
	}
	return fmt.Sprintf("<span class=\"%s\">%.*f%%</span>", class, decimalPlaces, stat)
}
