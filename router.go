package pong

import (
	"fmt"
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
}

func newRouter(pong *Pong) *Router {
	return &Router{
		pong:          pong,
		subRoutersMap: make(map[string]*Router),
		subHandlesMap: make(map[subHandlesMapKey]HandleFunc),
	}
}

func (r *Router) registerRouter(steps []string) *Router {
	parent, child := r, r
	for _, step := range steps {
		if len(step) > 0 && step[0] == ':' {
			child = parent.subRoutersMap[":"]
			if child == nil {
				for k, _ := range parent.subRoutersMap {
					fmt.Errorf("(%s) conflict (%s)\n", step, k+parent.paramName)
				}
				child = newRouter(parent.pong)
				parent.paramName = step[1:]
				parent.subRoutersMap[":"] = child
			}
		} else {
			child = parent.subRoutersMap[step]
			if child == nil {
				for k, _ := range parent.subRoutersMap {
					if k == step || k == ":" {
						fmt.Errorf("(%s) conflict (%s)\n", step, k+parent.paramName)
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

func (r *Router) registerHandle(step string, method string, handle HandleFunc) {
	if len(step) > 0 && step[0] == ':' {
		for k, _ := range r.subHandlesMap {
			if k.method == method && k.path == ":" {
				fmt.Errorf("(%s %s) conflict (%s %s)\n", step, method, k.path+r.paramName, k.method)
			}
		}
		r.paramName = step[1:]
		r.subHandlesMap[subHandlesMapKey{":", method}] = handle
	} else {
		for k, _ := range r.subHandlesMap {
			if (step == k.path && method == k.method) || k.path == ":" {
				fmt.Errorf("(%s %s) conflict (%s %s)\n", step, method, k.path+r.paramName, k.method)
			}
		}
		r.subHandlesMap[subHandlesMapKey{step, method}] = handle
	}
}

func (r *Router) register(path string, method string, handle HandleFunc) {
	steps := splitPath(path)
	stepsLength := len(steps)
	if stepsLength == 1 {
		r.registerHandle(steps[0], method, handle)
	} else {
		lastIndex := stepsLength - 1
		router := r.registerRouter(steps[0:lastIndex])
		lastStep := steps[lastIndex]
		router.registerHandle(lastStep, method, handle)
	}
}

// add a Middleware to this router
// this Middleware list will execute in order before execute the handle you provide to response
// all of this router's sub router will also execute this Middleware list,parent's Middleware list first child's Middleware list later
func (r *Router) Middleware(handles ...HandleFunc) {
	r.middlewareList = append(r.middlewareList, handles...)
}

// Add a sub router to this router
func (r *Router) Router(path string) *Router {
	steps := splitPath(path)
	return r.registerRouter(steps)
}

// register an path to handle HTTP Delete request
func (r *Router) Delete(path string, handle HandleFunc) {
	r.register(path, http.MethodDelete, handle)
}

// register an path to handle HTTP Get request
func (r *Router) Get(path string, handle HandleFunc) {
	r.register(path, http.MethodGet, handle)
}

// register an path to handle HTTP Head request
func (r *Router) Head(path string, handle HandleFunc) {
	r.register(path, http.MethodHead, handle)
}

// register an path to handle HTTP Options request
func (r *Router) Options(path string, handle HandleFunc) {
	r.register(path, http.MethodOptions, handle)
}

// register an path to handle HTTP Patch request
func (r *Router) Patch(path string, handle HandleFunc) {
	r.register(path, http.MethodPatch, handle)
}

// register an path to handle HTTP Post request
func (r *Router) Post(path string, handle HandleFunc) {
	r.register(path, http.MethodPost, handle)
}

// register an path to handle HTTP Put request
func (r *Router) Put(path string, handle HandleFunc) {
	r.register(path, http.MethodPut, handle)
}

// register an path to handle HTTP Trace request
func (r *Router) Trace(path string, handle HandleFunc) {
	r.register(path, http.MethodTrace, handle)
}

// register an path to handle any type HTTP request
// incloud "GET" "HEAD" "POST" "PUT" "PATCH" "DELETE" "CONNECT" "OPTIONS" "TRACE"
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
	nowStep := steps[0]
	isParamStep := false
	if len(r.paramName) > 0 {
		context.Request.requestParamMap[r.paramName] = nowStep
		isParamStep = true
	}
	if stepsLength == 1 {
		handleKey := subHandlesMapKey{
			path:   nowStep,
			method: context.Request.HTTPRequest.Method,
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
	} else {
		subRouter := r.subRoutersMap[nowStep]
		if subRouter == nil && isParamStep {
			subRouter = r.subRoutersMap[":"]
		}
		if subRouter != nil {
			subRouter.handle(steps[1:], context)
			return
		}
	}
	context.pong.NotFindHandle(context) //404
}
