/*
pong is a simple HTTP router for go.

Example:

    package main

    import (
    	"github.com/gwuhaolin/pong"
    	"net/http"
    	"log"
    )

    func main() {
    	po := pong.New()

    	// visitor http://127.0.0.1:3000/ping will see string "pong"
    	po.Root.Get("/ping", func(c *pong.Context) {
    		c.Response.String("pong")
    	})

    	// a sub router
    	sub := po.Root.Router("/sub")

    	// visitor http://127.0.0.1:3000/sub/pong will see JSON "{"name":"pong"}"
    	sub.Get("/:name", func(c *pong.Context) {
    		m := map[string]string{
    			"name":c.Request.Param("name"),
    		}
    		c.Response.JSON(m)
    	})

    	// Run Server Listen on 127.0.0.1:3000
    	log.Println(http.ListenAndServe(":3000", po))
    }

Learn more at https://github.com/gwuhaolin/pong
*/
package pong

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var (
// SessionId's Cookies name store in browser
	SessionCookiesName = "SESSIONID"
// this error will be return when use bind in request when bind data to struct fail
	ErrorTypeNotSupport = errors.New("type not support")
)

type (
// HandleFunc is a handle in Middleware list, like a machine production line to do some change
// used to read something from request and store by Context.Request
// make a response to client by Context.Response
	HandleFunc func(*Context)
	Pong       struct {
		htmlTemplate       *template.Template
		tailMiddlewareList []HandleFunc
		// Root router to path /
		Root               *Router
		// 404 not find handle
		// when pong's router can't find a handle to request' URL,pong will use NotFindHandle to handle this request
		// default is response with code 404, and string page not find
		NotFindHandle      HandleFunc
		// when send response to client cause error happen, pong will use HTTPErrorHandle to handle this request
		// default is response with code 500, and string inter server error
		HTTPErrorHandle    func(error, *Context)
		// SessionManager used to store and update value in session when pong has EnableSession
		// default SessionManager store data in memory
		SessionManager     SessionManager
	}
)

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	return strings.Split(path, "/")
}

// make a pong instance and return is pointer.
func New() *Pong {
	pong := &Pong{
		NotFindHandle: func(c *Context) {
			http.NotFound(c.Response.HTTPResponseWriter, c.Request.HTTPRequest)
		},
		HTTPErrorHandle: func(err error, c *Context) {
			c.Response.StatusCode = http.StatusInternalServerError
			c.Response.String(err.Error())
		},
	}
	pong.Root = newRouter(pong)
	return pong
}

// http.Server's ListenAndServe Handler
func (pong *Pong) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	steps := splitPath(request.URL.Path)
	context := newContext(pong, writer, request)
	pong.Root.handle(steps, context)
}

// load HTML template files whit glob
// if you will use render in response,you must call LoadTemplateGlob first to load template files.
// LoadTemplateGlob creates a new Template and parses the template definitions from the
// files identified by the pattern, which must match at least one file. The
// returned template will have the (base) name and (parsed) contents of the
// first file matched by the pattern. LoadTemplateGlob is equivalent to calling
// ParseFiles with the list of files matched by the pattern.
func (pong *Pong) LoadTemplateGlob(path string) {
	htmlTemplate, err := template.ParseGlob(path)
	if err != nil {
		fmt.Errorf("pong:%v", err)
	}
	pong.htmlTemplate = htmlTemplate
}

// add a middleware in the process's tail.
// which will execute before response data to client and after all of the other middleware register in router
// if you add more than one middlewares,this middlewares will execute in order
func (pong *Pong) TailMiddleware(middlewareList ...HandleFunc) {
	pong.tailMiddlewareList = append(pong.tailMiddlewareList, middlewareList...)
}
