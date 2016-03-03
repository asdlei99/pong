package pong

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"github.com/gwuhaolin/pong/_test"
)

func httpGetAssert(path string, responseStr string, t *testing.T) {
	res, err := http.Get(path)
	if err != nil {
		t.Error(err)
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	result := string(bs)
	if result != responseStr {
		t.Error(result, responseStr)
	}
}

func httpPostAssert(path string, contentType, bodyStr string, responseStr string, t *testing.T) {
	res, err := http.Post(path, contentType, strings.NewReader(responseStr))
	if err != nil {
		t.Error(err)
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	result := string(bs)
	if result != responseStr {
		t.Error(result, responseStr)
	}
}

func TestSplitPath(t *testing.T) {
	if len(splitPath("")) != 1 {
		t.Error()
	}
	if len(splitPath("/")) != 1 {
		t.Error()
	}
	if len(splitPath("//")) != 1 {
		t.Error()
	}
	if len(splitPath("/a/")) != 1 {
		t.Error()
	}
	if len(splitPath("a/b")) != 2 {
		t.Error()
	}
}

func TestRouter(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	root.Get("", func(c *Context) {
		c.Response.String("")
	})
	root.Get("/", func(c *Context) {
		c.Response.String("/")
	})
	defer httpGetAssert(baseURL + "", "/", t)
	defer httpGetAssert(baseURL + "/", "/", t)
	root.Post("/", func(c *Context) {
		c.Response.String("POST /")
	})
	defer httpPostAssert(baseURL + "/", applicationForm, "", "POST /", t)

	root.Get("/hi", func(c *Context) {
		c.Response.String("/hi")
	})
	defer httpGetAssert(baseURL + "/hi", "/hi", t)

	root.Post("/hi", func(c *Context) {
		c.Response.String("POST /hi")
	})
	defer httpPostAssert(baseURL + "/hi", applicationForm, "", "POST /hi", t)

	root.Get("/query", func(c *Context) {
		c.Response.String("/query?name=" + c.Request.Query("name"))
	})
	defer httpGetAssert(baseURL + "/query", "/query?name=", t)
	defer httpGetAssert(baseURL + "/query?name=吴浩麟", "/query?name=吴浩麟", t)

	root.Any("/any", func(c *Context) {
		c.Response.String("/any")
	})
	defer httpGetAssert(baseURL + "/any", "/any", t)
	defer httpPostAssert(baseURL + "/any", applicationForm, "", "/any", t)

	root.Get("/:param", func(c *Context) {
		c.Response.String("/:" + c.Request.Param("param"))
	})
	defer httpGetAssert(baseURL + "/a", "/:a", t)

	root.Get("/param/:id", func(c *Context) {
		c.Response.String("/param/:" + c.Request.Param("id"))
	})
	defer httpGetAssert(baseURL + "/param/123", "/param/:123", t)

	root.Get("/param/a/:id", func(c *Context) {
		c.Response.String("/param/a/:" + c.Request.Param("id"))
	})
	defer httpGetAssert(baseURL + "/param/a/123", "/param/a/:123", t)

	root.Get("/user/:id/update/:data", func(c *Context) {
		c.Response.String("/user/" + c.Request.Param("id") + "/update/" + c.Request.Param("data"))
	})
	defer httpGetAssert(baseURL + "/user/123/update/{age:24}", "/user/123/update/{age:24}", t)

	root.Get("/note/:id/update/:data", func(c *Context) {
		c.Response.String("/note/" + c.Request.Param("id") + "/update/" + c.Request.Param("data"))
	})
	defer httpGetAssert(baseURL + "/note/123/update/{age:24}", "/note/123/update/{age:24}", t)

	root.Get("/note/:id/remove", func(c *Context) {
		c.Response.String("/note/" + c.Request.Param("id") + "/remove")
	})
	defer httpGetAssert(baseURL + "/note/123/remove", "/note/123/remove", t)

	sub := root.Router("/sub")

	sub.Get("/hi", func(c *Context) {
		c.Response.String("/sub/hi")
	})
	defer httpGetAssert(baseURL + "/sub/hi", "/sub/hi", t)

	sub.Get("/note/:param", func(c *Context) {
		c.Response.String("/sub/note/:" + c.Request.Param("param"))
	})
	defer httpGetAssert(baseURL + "/sub/note/a", "/sub/note/:a", t)

	sub2 := sub.Router("/sub2")
	sub2.Get("", func(c *Context) {
		c.Response.String("/sub/sub2")
	})
	defer httpGetAssert(baseURL + "/sub/sub2", _test_util.NotFindString, t)

	sub2.Get("/", func(c *Context) {
		c.Response.String("/sub/sub2/")
	})
	defer httpGetAssert(baseURL + "/sub/sub2/", _test_util.NotFindString, t)

	sub2.Get("/hi", func(c *Context) {
		c.Response.String("/sub/sub2/hi")
	})
	defer httpGetAssert(baseURL + "/sub/sub2/hi", "/sub/sub2/hi", t)

	sub2.Post("/:param/hi", func(c *Context) {
		c.Response.String("POST /sub/sub2/:" + c.Request.Param("param"))
	})
	defer httpPostAssert(baseURL + "/sub/sub2/中文/hi", applicationForm, "", "POST /sub/sub2/:中文", t)

	sub2.Any("/user/any", func(c *Context) {
		c.Response.String("/sub/sub2/user/any")
	})
	defer httpGetAssert(baseURL + "/sub/sub2/user/any", "/sub/sub2/user/any", t)
	defer httpPostAssert(baseURL + "/sub/sub2/user/any", applicationForm, "", "/sub/sub2/user/any", t)
}

// /:name conflict with /path, but use /path first
func TestRouterConflict_Handle_PathOverParam(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	root.Get("/path", func(c *Context) {
		c.Response.String("path")
	})
	root.Get("/:name", func(c *Context) {
		c.Response.String(c.Request.Param("name"))
	})
	defer httpGetAssert(baseURL + "/path", "path", t)
	defer httpGetAssert(baseURL + "/hal", "hal", t)
}

// /:name conflict with /path, but use /path first
func TestRouterConflict_Handle_ParamOverParam(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	root.Get("/:path", func(c *Context) {
		c.Response.String("path=" + c.Request.Param("path"))
	})
	root.Get("/:name", func(c *Context) {
		c.Response.String("name=" + c.Request.Param("name"))
	})
	defer httpGetAssert(baseURL + "/abc", "name=abc", t)
}

func TestRouterMW(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	c := root.Router("a/b/c")
	c.Get("hi", func(c *Context) {
		c.Response.String("c")
	})
	a := root.Router("a")
	a.Middleware(func(c *Context) {
		c.Response.String("a")
	})
	b := root.Router("a/b")
	b.Middleware(func(c *Context) {
		c.Response.String("b")
	})
	defer httpGetAssert(baseURL + "/a/b/c/hi", "abc", t)
}

func TestParamInRouter(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	c := root.Router("a/:b/c")
	b := root.Router("a/:")
	b.Middleware(func(c *Context) {
		c.Response.String("b")
	})
	c.Get("hi", func(c *Context) {
		c.Response.String(c.Request.Param("b"))
	})
	defer httpGetAssert(baseURL + "/a/hal/c/hi", "bhal", t)
}

func TestRouterConflict(t *testing.T) {
	po, _ := runPong()
	root := po.Root
	root.Router("/a")
	root.Router("/a/:b")
	root.Router("/a/:b/c")
	root.Router("/a/:b/c")
	root.Router("/:a/")
	root.Get("/:a", po.NotFindHandle)
	root.Get("/hi", po.NotFindHandle)
	root.Get("/hi", po.NotFindHandle)
	root.Get("/:b", po.NotFindHandle)
}

func TestHead(t *testing.T) {
	po, baseURL := runPong()
	po.Root.Head("/", func(c *Context) {
		t.Log(`TestDelete`)
	})
	defer func() {
		_, err := http.Head(baseURL)
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestDelete(t *testing.T) {
	po, baseURL := runPong()
	po.Root.Delete("/", func(c *Context) {
		t.Log(`TestDelete`)
	})
	defer func() {
		client := http.Client{}
		url, _ := url.Parse(baseURL)
		_, err := client.Do(&http.Request{
			Method: http.MethodDelete,
			URL:    url,
		})
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestOptions(t *testing.T) {
	po, baseURL := runPong()
	po.Root.Options("/", func(c *Context) {
		t.Log(`TestOptions`)
	})
	defer func() {
		client := http.Client{}
		url, _ := url.Parse(baseURL)
		client.Do(&http.Request{
			Method: http.MethodOptions,
			URL:    url,
		})
	}()
}

func TestPatch(t *testing.T) {
	po, baseURL := runPong()
	po.Root.Patch("/", func(c *Context) {
		t.Log(`TestPatch`)
	})
	defer func() {
		client := http.Client{}
		url, _ := url.Parse(baseURL)
		client.Do(&http.Request{
			Method: http.MethodPatch,
			URL:    url,
		})
	}()
}

func TestPut(t *testing.T) {
	po, baseURL := runPong()
	po.Root.Put("/", func(c *Context) {
		t.Log(`TestPut`)
	})
	defer func() {
		client := http.Client{}
		url, _ := url.Parse(baseURL)
		client.Do(&http.Request{
			Method: http.MethodPut,
			URL:    url,
		})
	}()
}

func TestTrace(t *testing.T) {
	po, baseURL := runPong()
	po.Root.Trace("/", func(c *Context) {
		t.Log(`TestTrace`)
	})
	defer func() {
		client := http.Client{}
		url, _ := url.Parse(baseURL)
		client.Do(&http.Request{
			Method: http.MethodTrace,
			URL:    url,
		})
	}()
}
