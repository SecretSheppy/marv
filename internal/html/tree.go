package html

import (
	"os"
	"path"
	"sort"
	"strings"
)

type NodeType int

const (
	Directory NodeType = iota
	File
)

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
