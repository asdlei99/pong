package pong

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	SessionCookiesName  = "SESSIONID"
	ErrorTypeNotSupport = errors.New("type not support")
)

type (
	HandleFunc func(*Context)
	Pong       struct {
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
	} else {
		return []string{}
	}
}

// make a pong instance and return is pointer.
func New() *Pong {
	pong := &Pong{
		Logger: log.New(os.Stdout, "Pong:", log.LstdFlags),
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

//ignore this method
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
		pong.Logger.Println(err)
	}
	pong.htmlTemplate = htmlTemplate
}

//add a middleware in the process's tail.
//which will execute before response data to client and after all of the other middleware register in router
//if you add more than one middlewares,this middlewares will execute in order
func (pong *Pong) TailMiddleware(middlewareList ...HandleFunc) {
	pong.tailMiddlewareList = append(pong.tailMiddlewareList, middlewareList...)
}
