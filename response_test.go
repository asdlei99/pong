package pong

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"github.com/gwuhaolin/pong/_test"
)

func TestHeader(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	root.Get("/hi", func(c *Context) {
		c.Response.Header("X-name", "mine header")
		c.Response.String("TestHeader")
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		}
		header := res.Header.Get("X-name")
		if header != "mine header" {
			t.Error(header)
		}
		t.Log(`TestHeader`)
	}()
}

func TestCookie(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	root.Get("/hi", func(c *Context) {
		c.Response.Cookie(&http.Cookie{Name: "id", Value: "123"})
		c.Response.String("")
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		} else {
			cookies := res.Cookies()
			if len(cookies) != 1 || cookies[0].Name != "id" || cookies[0].Value != "123" {
				t.Error(cookies)
			}
		}
		t.Log(`TestCookie`)
	}()
}

func TestJSON(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := _test_util.TestUser{
		Name:  "吴浩麟",
		Age:   23,
		Money: 123.456,
		Alive: true,
		Notes: []_test_util.TestNote{
			{Text: "明天去放风筝"},
			{Text: "今天我们去逛宜家啦"},
		},
	}
	root.Get("/hi", func(c *Context) {
		c.Response.JSON(user)
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		} else {
			res_bs, _ := ioutil.ReadAll(res.Body)
			res.Body.Close()
			real_bs, _ := json.Marshal(user)
			if !reflect.DeepEqual(res_bs, real_bs) {
				t.Error(string(res_bs))
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != applicationJSONCharsetUTF8 {
				t.Error(ct)
			}
			t.Log(`TestJSON`)
		}
	}()
}

func TestJSONP(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := _test_util.TestUser{
		Name:  "吴浩麟",
		Age:   23,
		Money: 123.456,
		Alive: true,
		Notes: []_test_util.TestNote{
			{Text: "明天去放风筝"},
			{Text: "今天我们去逛宜家啦"},
		},
	}
	root.Get("/hi", func(c *Context) {
		cb := c.Request.Query("CALLBACK")
		c.Response.JSONP(user, cb)
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi?CALLBACK=hello")
		if err != nil {
			t.Error(err)
		} else {
			res_bs, _ := ioutil.ReadAll(res.Body)
			res_str := string(res_bs)
			res.Body.Close()
			real_bs, _ := json.Marshal(user)
			real_str := string(real_bs)
			if res_str != ("hello(" + real_str + ")") {
				t.Error(string(res_str))
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != applicationJavaScriptCharsetUTF8 {
				t.Error(ct)
			}
			t.Log(`TestJSONP`)
		}
	}()
}

func TestXML(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := _test_util.TestUser{
		Name:  "吴浩麟",
		Age:   23,
		Money: 123.456,
		Alive: true,
		Notes: []_test_util.TestNote{
			{Text: "明天去放风筝"},
			{Text: "今天我们去逛宜家啦"},
		},
	}
	root.Get("/hi", func(c *Context) {
		c.Response.XML(user)
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		} else {
			res_bs, _ := ioutil.ReadAll(res.Body)
			res.Body.Close()
			real_bs, _ := xml.Marshal(user)
			if !reflect.DeepEqual(res_bs, real_bs) {
				t.Error(string(res_bs))
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != applicationXMLCharsetUTF8 {
				t.Error(ct)
			}
			t.Log(`TestXML`)
		}
	}()
}

func TestFile(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	root.Get("/file", func(c *Context) {
		c.Response.File("_test/html/index.html")
	})
	defer func() {
		res, err := http.Get(baseURL + "/file")
		if err != nil {
			t.Error(err)
		} else {
			bs, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err == nil {
				str := string(bs)
				if str != "<h1>index.html</h1><b>{{.}}</b>" {
					t.Error(str)
				}
			}
			t.Log(`TestFile`)
		}
	}()
}

func TestString(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	bodyStr := "hello,乓"
	root.Get("/hi", func(c *Context) {
		c.Response.String(bodyStr)
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		} else {
			bs, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err == nil {
				str := string(bs)
				if str != bodyStr {
					t.Error(str)
				}
				ct := res.Header.Get(httpHeaderContentType)
				if ct != textPlainCharsetUTF8 {
					t.Error(ct)
				}
			}
			t.Log(`TestString`)
		}
	}()
}

func TestHTML(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	bodyHTML := "<h1>hello,乓</>"
	root.Get("/hi", func(c *Context) {
		c.Response.HTML(bodyHTML)
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		} else {
			bs, _ := ioutil.ReadAll(res.Body)
			res.Body.Close()
			html := string(bs)
			if html != bodyHTML {
				t.Error(html)
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != textHTMLCharsetUTF8 {
				t.Error(ct)
			}
			t.Log(`TestHTML`)
		}
	}()
}

func TestRender(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	po.LoadTemplateGlob("_test/html/*.html")
	root.Get("/render/:name", func(c *Context) {
		name := c.Request.Param("name")
		c.Response.Render(name, "中文")
	})
	defer func() {
		res, err := http.Get(baseURL + "/render/index.html")
		if err != nil {
			t.Error(err)
		} else {
			bs, _ := ioutil.ReadAll(res.Body)
			res.Body.Close()
			html := string(bs)
			if html != "<h1>index.html</h1><b>中文</b>" {
				t.Error(html)
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != textHTMLCharsetUTF8 {
				t.Error(ct)
			}
			t.Log(`TestRender`)
		}
	}()
	defer func() {
		res, err := http.Get(baseURL + "/render/footer.html")
		if err != nil {
			t.Error(err)
		} else {
			bs, _ := ioutil.ReadAll(res.Body)
			res.Body.Close()
			html := string(bs)
			if html != "<h1>footer.html</h1><b>中文</b>" {
				t.Error(html)
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != textHTMLCharsetUTF8 {
				t.Error(ct)
			}
			t.Log(`TestRender`)
		}
	}()
}

func TestRedirect(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	root.Get("/redirect", func(c *Context) {
		c.Response.Redirect("/")
	})
	defer func() {
		res, err := http.Get(baseURL + "/redirect")
		if err != nil {
			t.Error(err)
		} else {
			bs, _ := ioutil.ReadAll(res.Body)
			res.Body.Close()
			html := string(bs)
			if !strings.Contains(html, `a href="/"`) {
				t.Error(html)
			}
			t.Log(`TestRedirect`)
		}
	}()
}
