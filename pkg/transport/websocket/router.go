package websocket

import "path/filepath"

func NewRouter() *Router {
	return &Router{routeMap: map[string]HandleFunc{}}
}

type Router struct {
	routeMap map[string]HandleFunc
}

func (r *Router) Group(subAction string) *GroupRouter {
	return &GroupRouter{
		router:     r,
		baseAction: filepath.Join("", subAction),
	}
}

func (r *Router) Accept(action string, hf HandleFunc) {
	r.routeMap[action] = hf
}

func (r *Router) find(action string) (hf HandleFunc, ok bool) {
	hf, ok = r.routeMap[action]

	return hf, ok
}

type GroupRouter struct {
	router *Router

	baseAction string
}

func (r *GroupRouter) Group(subAction string) *GroupRouter {
	return &GroupRouter{
		router:     r.router,
		baseAction: filepath.Join(r.baseAction, subAction),
	}
}

func (r *GroupRouter) Accept(action string, hf HandleFunc) {
	r.router.Accept(filepath.Join(r.baseAction, action), hf)
}
