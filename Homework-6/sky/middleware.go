package sky

import (
	"log"
	"net/http"
)

func Logger() HandlerFunc {
	return func(ctx *Context) {
		ctx.code = 200
		ctx.Next()
		debugPrint("[%s]    %-100s%3d", ctx.Req.Method, ctx.Req.URL.Path, ctx.code)
	}
}

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				ctx.String(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		ctx.Next()
	}
}
