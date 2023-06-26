package cill

import "path"

/*
	RouterGroup 用于给路由分组
	提供功能
		（1）创建分组
		（2）使用中间件
		（3）创建POST GET DELETE PUT ... 路由handle映射
*/
type RouterGroup struct {
	pattern string
	middlewares []HandlerFunc
	childGroups []*RouterGroup
	engine *Engine
}


func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *RouterGroup) Group(comp string) *RouterGroup {
	newGroup := &RouterGroup{
		pattern: path.Join(g.pattern, comp),
		engine: g.engine,
	}
	g.childGroups = append(g.childGroups, newGroup)
	return newGroup
}

func (g *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := path.Join(g.pattern, comp)
	g.engine.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(comp string, handler HandlerFunc) {
	g.addRoute("GET", comp, handler)
}

func (g *RouterGroup) POST(comp string, handler HandlerFunc) {
	g.addRoute("POST", comp, handler)
}