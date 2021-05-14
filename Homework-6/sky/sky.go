package sky

import (
	"net/http"
	"path"
	"strings"
)

const (
	default404body = "404 Not Found"
)

type Engine struct {
	*RouterGroup
	methodTrees methodTrees
}

func Default() *Engine {
	e := new(Engine)
	e.RouterGroup = newRouterGroup()
	e.RouterGroup.engine = e
	e.RouterGroup.isRoot = true
	e.RouterGroup.prefix = "/"
	e.Use(Logger(), Recovery())
	return e
}

func (e *Engine) handle(method string, url string) (handlers HandlerFuncs, param map[string]string) {
	handlers, originUrl := e.RouterGroup.handle(method, url)
	if handlers == nil {
		return nil, nil
	}

	param = make(map[string]string)
	originLice := strings.Split(originUrl, "/")[1:]
	addrLice := strings.Split(url, "/")[1:]
	for index, v := range originLice {
		if len(v) != 0 {
			if v[0] == ':' {
				param[v[1:]] = addrLice[index]
			} else if v[0] == '*' {
				param[v[1:]] = path.Join(addrLice[index:]...)
			}
		}
	}

	return
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	ctx := newContext(writer, req)
	ctx.handlers, ctx.param = e.handle(ctx.Req.Method, ctx.Req.URL.Path)
	if len(ctx.handlers) != 0 {
		ctx.Next()
	} else {
		if isDebugging() {
			debugPrint("[%s]    %-100s%3d", ctx.Req.Method, ctx.Req.URL.Path, 404)
		}
		ctx.String(404, default404body)
	}
}

func (e *Engine) Run(addr ...string) (err error) {
	defer debugPrintError(err)

	var port string
	switch len(addr) {
	case 0:
		port = ":8080"
		debugPrint("Listening address is not set. It is set to \":8080\"")
	case 1:
		port = addr[0]
		debugPrint("The serve is set to \"" + port + "\"")
	default:
		panic("Listening address is too many")
	}

	err = http.ListenAndServe(port, e)
	return err
}
