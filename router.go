package pong

import (
	"net/http"
)

type subHandlesMapKey struct {
	path   string
	method string
}

type Router struct {
	pong           *Pong
	paramName      string
	middlewareList []HandleFunc
	subRoutersMap  map[string]*Router
	subHandlesMap  map[subHandlesMapKey]HandleFunc
	indexHandleMap map[string]HandleFunc
}

func newRouter(pong *Pong) *Router {
	return &Router{
		pong:pong,
		subRoutersMap:make(map[string]*Router),
		subHandlesMap:make(map[subHandlesMapKey]HandleFunc),
		indexHandleMap:make(map[string]HandleFunc),
	}
}

func (r *Router)registerRouter(steps []string) *Router {
	parent, child := r, r
	for _, step := range steps {
		if step[0] == ':' {
			child = parent.subRoutersMap[":"]
			if child == nil {
				for k, _ := range parent.subRoutersMap {
					r.pong.Logger.Printf("(%s) conflict (%s)\n", step, k + parent.paramName)
				}
				child = newRouter(parent.pong)
				parent.paramName = step[1:]
				parent.subRoutersMap[":"] = child
			}
		}else {
			child = parent.subRoutersMap[step]
			if child == nil {
				for k, _ := range parent.subRoutersMap {
					if k == step || k == ":" {
						parent.pong.Logger.Printf("(%s) conflict (%s)\n", step, k + parent.paramName)
					}
				}
				child = newRouter(parent.pong)
				parent.subRoutersMap[step] = child
			}
		}
		parent = child
	}
	return child
}

func (r *Router)registerHandle(step string, method string, handle HandleFunc) {
	if step[0] == ':' {
		for k, _ := range r.subHandlesMap {
			if k.method == method && k.path == ":" {
				r.pong.Logger.Printf("(%s %s) conflict (%s %s)\n", step, method, k.path + r.paramName, k.method)
			}
		}
		r.paramName = step[1:]
		r.subHandlesMap[subHandlesMapKey{":", method}] = handle
	}else {
		for k, _ := range r.subHandlesMap {
			if (step == k.path && method == k.method) || k.path == ":" {
				r.pong.Logger.Printf("(%s %s) conflict (%s %s)\n", step, method, k.path + r.paramName, k.method)
			}
		}
		r.subHandlesMap[subHandlesMapKey{step, method}] = handle
	}
}

func (r *Router)register(path string, method string, handle HandleFunc) {
	steps := splitPath(path)
	stepsLength := len(steps)
	switch stepsLength {
	case 0:
		r.indexHandleMap[method] = handle
	case 1:
		r.registerHandle(steps[0], method, handle)
	default:
		lastIndex := stepsLength - 1
		router := r.registerRouter(steps[0:lastIndex])
		lastStep := steps[lastIndex]
		router.registerHandle(lastStep, method, handle)
	}
}

func (r *Router)Middleware(handles ...HandleFunc) {
	r.middlewareList = append(r.middlewareList, handles...)
}

func (r *Router)Router(path string) *Router {
	steps := splitPath(path)
	return r.registerRouter(steps)
}

func (r *Router) Delete(path string, handle HandleFunc) {
	r.register(path, http.MethodDelete, handle)
}

func (r *Router) Get(path string, handle HandleFunc) {
	r.register(path, http.MethodGet, handle)
}

func (r *Router) Head(path string, handle HandleFunc) {
	r.register(path, http.MethodHead, handle)
}

func (r *Router) Options(path string, handle HandleFunc) {
	r.register(path, http.MethodOptions, handle)
}

func (r *Router) Patch(path string, handle HandleFunc) {
	r.register(path, http.MethodPatch, handle)
}

func (r *Router) Post(path string, handle HandleFunc) {
	r.register(path, http.MethodPost, handle)
}

func (r *Router) Put(path string, handle HandleFunc) {
	r.register(path, http.MethodPut, handle)
}

func (r *Router) Trace(path string, handle HandleFunc) {
	r.register(path, http.MethodTrace, handle)
}

func (r *Router) Any(path string, handle HandleFunc) {
	r.register(path, http.MethodDelete, handle)
	r.register(path, http.MethodGet, handle)
	r.register(path, http.MethodHead, handle)
	r.register(path, http.MethodOptions, handle)
	r.register(path, http.MethodPatch, handle)
	r.register(path, http.MethodPost, handle)
	r.register(path, http.MethodPut, handle)
	r.register(path, http.MethodTrace, handle)
}

func (r *Router) handle(steps []string, context *Context) {
	for _, handle := range r.middlewareList {
		handle(context)
	}
	stepsLength := len(steps)
	if stepsLength == 0 {
		indexHandle := r.indexHandleMap[context.Request.HTTPRequest.Method]
		if indexHandle != nil {
			indexHandle(context)
			return
		}
	}else {
		nowStep := steps[0]
		isParamStep := false
		if len(r.paramName) > 0 {
			context.Request.requestParamMap[r.paramName] = nowStep
			isParamStep = true
		}
		if stepsLength == 1 {
			handleKey := subHandlesMapKey{
				path:nowStep,
				method:context.Request.HTTPRequest.Method,
			}
			handle := r.subHandlesMap[handleKey]
			if handle == nil && isParamStep {
				handleKey.path = ":"
				handle = r.subHandlesMap[handleKey]
			}
			if handle != nil {
				handle(context)
				return
			}
		}else {
			subRouter := r.subRoutersMap[nowStep]
			if subRouter == nil && isParamStep {
				subRouter = r.subRoutersMap[":"]
			}
			if subRouter != nil {
				subRouter.handle(steps[1:], context)
				return
			}
		}
	}
	context.pong.NotFindHandle(context)//404
}