package pong

import (
	"net/http"
	"log"
	"os"
	"strings"
	"errors"
	"html/template"
)

var (
	SessionCookiesName = "SESSIONID"
	ErrorTypeNotSupport = errors.New("type not support")
)

type(
	HandleFunc func(*Context)
	Pong struct {
		htmlTemplate       *template.Template
		tailMiddlewareList []HandleFunc
		Root               *Router
		Logger             *log.Logger
		NotFindHandle      HandleFunc
		HTTPErrorHandle    func(error, *Context)
		SessionManager     SessionManager
	}
)

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if len(path) > 0 {
		return strings.Split(path, "/")
	}else {
		return []string{}
	}
}

func New() *Pong {
	pong := &Pong{
		Logger:log.New(os.Stdout, "Pong:", log.LstdFlags),
		NotFindHandle:func(c *Context) {
			http.NotFound(c.Response.HTTPResponseWriter, c.Request.HTTPRequest)
		},
		HTTPErrorHandle:func(err error, c *Context) {
			c.Response.StatusCode = http.StatusInternalServerError
			c.Response.String(err.Error())
		},
	}
	pong.Root = newRouter(pong)
	return pong
}

func (pong *Pong)ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	steps := splitPath(request.URL.Path)
	context := newContext(pong, writer, request)
	pong.Root.handle(steps, context)
}

func (pong *Pong)LoadTemplateGlob(path string) {
	htmlTemplate, err := template.ParseGlob(path)
	if err != nil {
		pong.Logger.Println(err)
	}
	pong.htmlTemplate = htmlTemplate
}

func (pong *Pong)TailMiddleware(middlewareList ...HandleFunc) {
	pong.tailMiddlewareList = append(pong.tailMiddlewareList, middlewareList...)
}