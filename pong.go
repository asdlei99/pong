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
	ErrorContentTypeNotSupport = errors.New("http content type not support")
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
		SessionManager:&memorySessionManager{
			store:make(map[string]map[string]interface{}),
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
		log.Fatal(err)
	}
	pong.htmlTemplate = htmlTemplate
}

func (pong *Pong)EnableSession() {
	if pong.SessionManager == nil {
		pong.SessionManager = &memorySessionManager{
			store:make(map[string]map[string]interface{}),
		}
	}
	pong.Root.Middleware(func(c *Context) {
		c.Session = &Session{
			pong:c.pong,
		}
		sCookie, err := c.Request.HTTPRequest.Cookie(SessionCookiesName)
		if err == nil {
			c.Session.id = sCookie.Value
			if c.pong.SessionManager.Has(c.Session.id) {
				v := c.pong.SessionManager.Read(c.Session.id)
				c.Session.store = v
			}else {
				goto noSessionID
			}
		}else {

		}
		noSessionID:{
			c.Session.id = c.pong.SessionManager.NewSession()
			c.Response.Cookie(&http.Cookie{
				HttpOnly:true,
				Name:SessionCookiesName,
				Value:c.Session.id,
			})
		}

	})
	pong.TailMiddleware(func(c *Context) {
		change := make(map[string]interface{})
		for _, name := range c.Session.hasChangeFlag {
			change[name] = c.Session.store[name]
		}
		c.pong.SessionManager.Write(c.Session.id, change)
	})
}

func (pong *Pong)TailMiddleware(middlewareList ...HandleFunc) {
	pong.tailMiddlewareList = append(pong.tailMiddlewareList, middlewareList...)
}