package pong

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

const (
	notFindString = "404 page not found\n"
)

var (
	listenPort int = 3000
)

func runPong() (pong *Pong, baseURL string) {
	pong = New()
	pong.Root.Middleware(func(c *Context) {
		req := c.Request.HTTPRequest
		fmt.Println(req.Method, req.Host, req.RequestURI)
	})
	serverHasRun := make(chan bool)
	go func() {
		listenAddr := "127.0.0.1:" + strconv.Itoa(listenPort)
		baseURL = "http://" + listenAddr
		serverHasRun <- true
		http.ListenAndServe(listenAddr, pong)
	}()
	<-serverHasRun
	listenPort++
	return
}

func TestHTTPErrorHandle(t *testing.T) {
	po, baseURL := runPong()
	errorStr := "宝宝,我错了"
	po.Root.Get("/json", func(c *Context) {
		c.Response.JSON(po.HTTPErrorHandle)
	})
	po.Root.Get("/jsonp", func(c *Context) {
		c.Response.JSONP(po.HTTPErrorHandle, errorStr)
	})
	po.Root.Get("/xml", func(c *Context) {
		c.Response.XML(po.HTTPErrorHandle)
	})
	po.Root.Get("/render", func(c *Context) {
		c.Response.Render(errorStr, nil)
	})
	defer func() {
		for _, path := range [...]string{"json", "jsonp", "xml"} {
			res, err := http.Get(baseURL + "/" + path)
			if err != nil {
				t.Error(err)
			}
			bs, _ := ioutil.ReadAll(res.Body)
			res_str := string(bs)
			if len(res_str) == 0 {
				t.Error(path, res_str)
			}
			if res.StatusCode != http.StatusInternalServerError {
				t.Error(path, res.StatusCode)
			}
		}
		http.Get(baseURL + "/render")
	}()
}

func TestLoadTemplateGlobError(t *testing.T) {
	po, baseURL := runPong()
	po.LoadTemplateGlob("/no/this/file/")
	po.LoadTemplateGlob("test_resource/html/*.html")
	po.Root.Get("/render", func(c *Context) {
		c.Response.Render("/no/this/file.html", nil)
	})
	defer func() {
		http.Get(baseURL + "/render")
	}()
}
