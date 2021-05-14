package sky

import (
	"path"
	"strings"
)

type node struct {
	group        *RouterGroup
	isRegistered bool
	handlers     HandlerFuncs
	fullPart     string
	part         string
	children     []*node
	exactMatch   bool
}

func newNode() *node {
	res := new(node)
	res.handlers = newHandlerFuncs()
	return res
}

func (n *node) matchExactChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

func (n *node) matchAllChild(part string) []*node {
	res := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || !child.exactMatch {
			res = append(res, child)
		}
	}
	return res
}

func (n *node) Insert(part string) *node {
	child := n.matchExactChild(part)
	if child == nil {
		child = newNode()
		child.part = part
		child.exactMatch = part[0] != '*' && part[0] != ':'
		child.fullPart = path.Join(n.fullPart, part)
		if n.children == nil {
			n.children = append(make([]*node, 0), child)
		} else {
			n.children = append(n.children, child)
		}
		return child
	} else {
		return child
	}
}

func (n *node) Search(addr string) *node {
	lice := strings.Split(addr, "/")[1:]
	return n.search(lice, 0)
}

func (n *node) search(lice []string, height int) *node {
	var res *node
	child := n.matchAllChild(lice[height])
	if len(child) == 0 {
		return nil
	}

	if height == len(lice)-1 || child[0].part[0] == '*' {
		if child[0].isRegistered {
			return child[0]
		} else {
			return nil
		}
	}

	for _, v := range child {
		res = v.search(lice, height+1)
		if res != nil {
			return res
		}
	}

	return nil
}
