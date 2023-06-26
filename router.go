package cill

import (
	"log"
	"net/http"
	"strings"
)

/*
	router 用于维护user定义的路由信息
	提供功能：
		(1) 添加路由信息
		(2）获取路由信息
		(3) 获取路由handle
		(4) 处理请求

*/
type router struct {
	root     map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		root:     make(map[string]*node),       // router 树用于匹配路由
		handlers: make(map[string]HandlerFunc), // url mapping handler
	}
}

// addRoute 添加route
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	_, ok := r.root[method]
	if !ok {
		r.root[method] = &node{}
	}
	r.root[method].insert(pattern, parts, 0)

	log.Printf("mount route [%s] %s\n", method, pattern)

	key := method + "-" + pattern
	r.handlers[key] = handler

	
}

func (r *router) getRoute(method, realPath string) (*node, map[string]string) {
	realPathParts := parsePattern(realPath)
	root, ok := r.root[method]
	if !ok {
		return nil, nil
	}
	n := root.Search(realPathParts, 0) 
	if n == nil {
		return nil, nil
	}

	// parse params
	params := make(map[string]string)
	patternParts := parsePattern(n.pattern)
	for k := range realPathParts {
		if strings.HasPrefix(patternParts[k], ":") {
			key := string(patternParts[k][1:])
			params[key] = realPathParts[k]
		}	
		if strings.HasPrefix(patternParts[k], "*") {
			key := string(patternParts[k][1:])
			params[key] = strings.Join(realPathParts[k:], "/")
		}
	}
	return n, params
}



func (r *router) getHandle(ctx *Context) HandlerFunc {
	n, params := r.getRoute(ctx.Method, ctx.Path)
	if n == nil {
		return func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 NOT FOUND. %s", ctx.Path)
		}
	}
	ctx.Params = params
	key := ctx.Method + "-" + n.pattern

	return r.handlers[key]
}


func (r *router) handle(ctx *Context) {
	ctx.handler = append(ctx.handler, r.getHandle(ctx))
	ctx.Next()
}

/*
	工具函数
*/

// parsePattern 解析url为parts
func parsePattern(pattern string) []string {
	res := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, v := range res {
		if v != "" {
			parts = append(parts, v)
		}
	}
	return parts
}