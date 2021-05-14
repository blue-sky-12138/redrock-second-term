package sky

import (
	"path"
	"strings"
)

const (
	GETMethod    = "GET"
	POSTMethod   = "POST"
	PUTMethod    = "PUT"
	DELETEMethod = "DELETE"
)

type HandlerFunc func(*Context)

type HandlerFuncs []HandlerFunc

func newHandlerFuncs() HandlerFuncs {
	return make(HandlerFuncs, 0)
}

type RouterGroup struct {
	isRoot   bool
	prefix   string
	engine   *Engine
	handlers HandlerFuncs
}

func newRouterGroup() *RouterGroup {
	res := new(RouterGroup)
	res.handlers = newHandlerFuncs()
	return res
}

func (r *RouterGroup) addRoute(method string, addr string, handlers ...HandlerFunc) {
	illegalPanic(r.isRoot && addr[0] != '/', "path must begin with '/'")
	illegalPanic(method == "", "method must not be empty")
	illegalPanic(len(handlers) == 0, "handler must have one at least")

	root := r.engine.methodTrees.getTree(method)
	if root == nil {
		root = newNode()
		root.children = make([]*node, 0)
		root.group = r
		root.fullPart = "/"
		r.engine.methodTrees = append(r.engine.methodTrees, methodTree{
			method: method,
			root:   root,
		})
	}

	if len(addr) != 0 && addr[0] != '/' {
		addr = "/" + addr
	}
	addr = path.Join(r.prefix, addr)
	lice := strings.Split(addr, "/")[1:]

	tempNode := root
	for index, part := range lice {
		if part != "" {
			illegalPanic(part[0] == '*' && index != len(lice)-1, "'*' must be the last of the addr")
		}
		tempNode = tempNode.Insert(part)
	}

	illegalPanic(tempNode.isRegistered, "you have two same router")
	tempNode.isRegistered = true
	tempNode.handlers = append(tempNode.handlers, handlers...)
	tempNode.group = r

	debugPrint("[" + method + "]    " + tempNode.fullPart)
}

func (r *RouterGroup) handle(method string, addr string) (HandlerFuncs, string) {
	search := r.getNode(method, addr)
	if search == nil {
		return nil, ""
	}

	res := search.group.handlers
	res = append(res, search.handlers...)

	return res, search.fullPart
}

func (r *RouterGroup) getNode(method string, addr string) *node {
	tree := r.engine.methodTrees.getTree(method)
	if tree == nil {
		return nil
	}

	return tree.Search(addr)
}

func (r *RouterGroup) Group(prefix string, handlers ...HandlerFunc) *RouterGroup {
	res := newRouterGroup()
	res.prefix = path.Join(r.prefix, prefix)
	res.engine = r.engine
	res.handlers = r.handlers
	res.handlers = append(res.handlers, handlers...)
	return res
}

func (r *RouterGroup) Use(handlers ...HandlerFunc) {
	r.handlers = append(r.handlers, handlers...)
}

func (r *RouterGroup) GET(addr string, handler ...HandlerFunc) {
	r.addRoute(GETMethod, addr, handler...)
}

func (r *RouterGroup) POST(addr string, handler ...HandlerFunc) {
	r.addRoute(POSTMethod, addr, handler...)
}

func (r *RouterGroup) PUT(addr string, handler ...HandlerFunc) {
	r.addRoute(PUTMethod, addr, handler...)
}

func (r *RouterGroup) DELETE(addr string, handler ...HandlerFunc) {
	r.addRoute(DELETEMethod, addr, handler...)
}
