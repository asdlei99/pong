package pong

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	po, baseURL := runPong()
	httpGetAssert := func(path string, responseStr string) {
		res, err := http.Get(baseURL + path)
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
		t.Log(`TestRouter`)
	}
	httpPostAssert := func(path string, contentType, bodyStr string, responseStr string) {
		res, err := http.Post(baseURL+path, contentType, strings.NewReader(responseStr))
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
		t.Log(`TestRouter`)
	}

	root := po.Root
	root.Get("", func(c *Context) {
		c.Response.String("")
	})
	root.Get("/", func(c *Context) {
		c.Response.String("/")
	})
	defer httpGetAssert("", "/")
	defer httpGetAssert("/", "/")
	root.Post("/", func(c *Context) {
		c.Response.String("POST /")
	})
	defer httpPostAssert("/", applicationForm, "", "POST /")

	root.Get("/hi", func(c *Context) {
		c.Response.String("/hi")
	})
	defer httpGetAssert("/hi", "/hi")

	root.Post("/hi", func(c *Context) {
		c.Response.String("POST /hi")
	})
	defer httpPostAssert("/hi", applicationForm, "", "POST /hi")

	root.Get("/query", func(c *Context) {
		c.Response.String("/query?name=" + c.Request.Query("name"))
	})
	defer httpGetAssert("/query", "/query?name=")
	defer httpGetAssert("/query?name=吴浩麟", "/query?name=吴浩麟")

	root.Any("/any", func(c *Context) {
		c.Response.String("/any")
	})
	defer httpGetAssert("/any", "/any")
	defer httpPostAssert("/any", applicationForm, "", "/any")

	root.Get("/:param", func(c *Context) {
		c.Response.String("/:" + c.Request.Param("param"))
	})
	defer httpGetAssert("/a", "/:a")

	root.Get("/param/:id", func(c *Context) {
		c.Response.String("/param/:" + c.Request.Param("id"))
	})
	defer httpGetAssert("/param/123", "/param/:123")

	root.Get("/param/a/:id", func(c *Context) {
		c.Response.String("/param/a/:" + c.Request.Param("id"))
	})
	defer httpGetAssert("/param/a/123", "/param/a/:123")

	root.Get("/user/:id/update/:data", func(c *Context) {
		c.Response.String("/user/" + c.Request.Param("id") + "/update/" + c.Request.Param("data"))
	})
	defer httpGetAssert("/user/123/update/{age:24}", "/user/123/update/{age:24}")

	root.Get("/note/:id/update/:data", func(c *Context) {
		c.Response.String("/note/" + c.Request.Param("id") + "/update/" + c.Request.Param("data"))
	})
	defer httpGetAssert("/note/123/update/{age:24}", "/note/123/update/{age:24}")

	root.Get("/note/:id/remove", func(c *Context) {
		c.Response.String("/note/" + c.Request.Param("id") + "/remove")
	})
	defer httpGetAssert("/note/123/remove", "/note/123/remove")

	sub := root.Router("/sub")

	sub.Get("/hi", func(c *Context) {
		c.Response.String("/sub/hi")
	})
	defer httpGetAssert("/sub/hi", "/sub/hi")

	sub.Get("/note/:param", func(c *Context) {
		c.Response.String("/sub/note/:" + c.Request.Param("param"))
	})
	defer httpGetAssert("/sub/note/a", "/sub/note/:a")

	sub2 := sub.Router("/sub2")
	sub2.Get("", func(c *Context) {
		c.Response.String("/sub/sub2")
	})
	defer httpGetAssert("/sub/sub2", notFindString)

	sub2.Get("/", func(c *Context) {
		c.Response.String("/sub/sub2/")
	})
	defer httpGetAssert("/sub/sub2/", notFindString)

	sub2.Get("/hi", func(c *Context) {
		c.Response.String("/sub/sub2/hi")
	})
	defer httpGetAssert("/sub/sub2/hi", "/sub/sub2/hi")

	sub2.Post("/:param/hi", func(c *Context) {
		c.Response.String("POST /sub/sub2/:" + c.Request.Param("param"))
	})
	defer httpPostAssert("/sub/sub2/中文/hi", applicationForm, "", "POST /sub/sub2/:中文")

	sub2.Any("/user/any", func(c *Context) {
		c.Response.String("/sub/sub2/user/any")
	})
	defer httpGetAssert("/sub/sub2/user/any", "/sub/sub2/user/any")
	defer httpPostAssert("/sub/sub2/user/any", applicationForm, "", "/sub/sub2/user/any")
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
