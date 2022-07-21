package cgin

import "net/http"

//Engine 实现Serve HTTP'接口使得能够替换默认的http.Server Handler
type Engine struct {
	//最基础的map路由表,每个http方法对应一个容纳Handler的map
	router map[string]map[string]CiginHandler
}

func NewEngine() *Engine {
	return &Engine{router: make(map[string]map[string]CiginHandler)}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//创建上下文
	ctx := NewContext(w, r)

	//路由表查询,对于一个web框架来讲,其路由表的结构,以及匹配的逻辑规则是最最重要的部分!
	if handlers, ok := e.router[r.Method]; ok {
		//对应的方法找到,查询第二层
		if handler, ok := handlers[r.URL.Path]; ok {
			//找到了 执行
			handler(ctx)
		} else {
			//没找到,返回404
			ctx.Json(http.StatusNotFound, "404 not found")
			return
		}

	} else {
		//没有相关的路由
		ctx.Json(http.StatusNotFound, "404 not found")
		return
	}

}

//restful风格的路由注册
func (e *Engine) Get(path string, handler CiginHandler) {

	if _, ok := e.router["GET"]; !ok {
		e.router["GET"] = make(map[string]CiginHandler)
	}
	e.router["GET"][path] = handler

}

/*
一个web框架的核心功能:
	1.高性能的路由表,包含路由表的数据结构,以及匹配的算法
		- 默认的是map路由表,两级map就能实现先根据方法,再根据uri的形式匹配,但无法支持动态路由
		- 为实现动态路由,需要使用前缀树,这就涉及到前缀树的设计,需要熟悉前缀树的增加 查找操作

	2.rest风格的路由注册,并且要同样前缀批量注册(即group分组)
	3.封装良好的上下文,快速的处理http请求
	4.合理的中间件支持,必须能够支持拓展,排除函数嵌套



*/
