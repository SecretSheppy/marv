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

func (p *PathNode) Render(buff *bytes.Buffer, fw fwlib.Framework, lang *languages.Language) {
	p.render(buff, fw, lang, 0, fmt.Sprintf("/%s/mutants", fw.Meta().Name))
}

func (p *PathNode) render(buff *bytes.Buffer, fw fwlib.Framework, lang *languages.Language, level int, accPath string) {
	currentPath := path.Join(accPath, p.Name)
	switch p.Type {
	case Directory:
		buff.WriteString("<div class=\"directory-wrapper\">")
		buff.WriteString(fmt.Sprintf("<div class=\"directory\" style=\"margin-left: %dpx\">%s</div>", level*5, p.Name))
		for _, child := range p.children {
			child.render(buff, fw, lang, level+1, currentPath)
		}
		buff.WriteString("</div>")
	case File:
		buff.WriteString(fmt.Sprintf("<a class=\"file\" style=\"margin-left: %dpx\" href=\"%s\">%s</a>", level*5, currentPath, p.Name))
	}
}
