package cill

import (
	"net/http"
	"strings"
)

type HandlerFunc func(ctx *Context)

/*
	Engine 统辖全部资源：routerGroup, router,
	Engine 本身也作为一个RouterGroup
*/
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func New() *Engine {
	e := &Engine{router: newRouter()}
	e.RouterGroup = &RouterGroup{engine: e}
	e.groups = append(e.groups, e.RouterGroup)
	return e
}

func Default() *Engine {
	e := New()
	e.Use(Logger(), Recover())
	return e
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := NewContext(w, r, e)

	// 添加 router_group middleware
	for _, g := range e.groups {
		if strings.HasPrefix(r.URL.Path, g.pattern) {
			ctx.handler = append(ctx.handler, g.middlewares...)
		}
	}

	// 开始处理
	e.router.handle(ctx)
}