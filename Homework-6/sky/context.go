package sky

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
)

type H map[string]interface{}

type Context struct {
	code     int
	index    int
	handlers HandlerFuncs
	param    map[string]string
	Writer   http.ResponseWriter
	Req      *http.Request
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		index:  -1,
		Writer: w,
		Req:    r,
	}
}

func (ctx *Context) Next() {
	ctx.index++
	l := len(ctx.handlers)
	for ; ctx.index < l; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

func (ctx *Context) Stop() {
	ctx.index = math.MaxInt8 / 2
}

func (ctx *Context) statusCode(code int) {
	ctx.code = code
	ctx.Writer.WriteHeader(code)
}

func (ctx *Context) SetHeader(key string, value string) {
	ctx.Writer.Header().Set(key, value)
}

func (ctx *Context) GetHeader(key string) string {
	return ctx.Req.Header.Get(key)
}

func (ctx *Context) Param(key string) string {
	return ctx.param[key]
}

func (ctx *Context) Params(h *H) {
	for key := range *h {
		(*h)[key] = ctx.param[key]
	}
}

func (ctx *Context) ParamAll() map[string]string {
	return ctx.param
}

func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

func (ctx *Context) PostForms(h *H) {
	for key := range *h {
		(*h)[key] = ctx.PostForm(key)
	}
}

func (ctx *Context) PostFormAll() H {
	values := ctx.Req.Form
	if values == nil {
		return H{}
	}

	res := make(H)
	for key, value := range values {
		if len(value) == 0 {
			res[key] = ""
		} else {
			res[key] = value[0]
		}
	}
	return res
}

func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

func (ctx *Context) Queries(h *H) {
	for key := range *h {
		(*h)[key] = ctx.Query(key)
	}
}

func (ctx *Context) QueryAll() H {
	values := ctx.Req.URL.Query()
	if values != nil {
		return H{}
	}

	res := make(H)
	for key, value := range values {
		if len(value) == 0 {
			res[key] = ""
		} else {
			res[key] = value[0]
		}
	}
	return res
}

func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetHeader("Content-Type", "application/json")
	ctx.statusCode(code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, "Server error", 500)
	}
}

func (ctx *Context) String(code int, format string, value ...interface{}) {
	ctx.SetHeader("Content-Type", "text/plain")
	ctx.statusCode(code)
	ctx.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

func (ctx *Context) HTML(code int, html string) {
	ctx.SetHeader("Content-Type", "text/html")
	ctx.statusCode(code)
	ctx.Writer.Write([]byte(html))
}

func (ctx *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: http.SameSiteDefaultMode,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (ctx *Context) GetCookie(key string) (string,error) {
	cookie, err := ctx.Req.Cookie(key)
	if err != nil {
		return "",err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (ctx *Context) GetCookies(h *H)(err error) {
	for key := range *h {
		(*h)[key],err = ctx.GetCookie(key)
		if err != nil {
			return
		}
	}
	return
}