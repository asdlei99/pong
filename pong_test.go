package pong

import (
	"io/ioutil"
	"net/http"
	"testing"
	"fmt"
	"strconv"
	"github.com/gwuhaolin/pong/_test"
)

func runPong() (po *Pong, baseURL string) {
	po = New()
	po.Root.Middleware(func(c *Context) {
		req := c.Request.HTTPRequest
		fmt.Println(req.Method, req.Host, req.RequestURI)
	})
	serverHasRun := make(chan bool)
	go func() {
		listenAddr := "127.0.0.1:" + strconv.Itoa(_test_util.ListenPort)
		baseURL = "http://" + listenAddr
		serverHasRun <- true
		http.ListenAndServe(listenAddr, po)
	}()
	<-serverHasRun
	_test_util.ListenPort++
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
	po.LoadTemplateGlob("_test/html/*.html")
	po.Root.Get("/render", func(c *Context) {
		c.Response.Render("/no/this/file.html", nil)
	})
	defer func() {
		http.Get(baseURL + "/render")
	}()
}
