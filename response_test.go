package pong

import (
	"testing"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"reflect"
	"encoding/xml"
)

func TestHeader(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	root := po.Root
	po.Root.Middleware(logRequest)
	root.Get("/hi", func(c *Context) {
		c.Response.Header("X-name", "mine header")
		c.Response.String("")
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
	}()
}

func TestCookie(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	root := po.Root
	po.Root.Middleware(logRequest)
	root.Get("/hi", func(c *Context) {
		c.Response.Cookie(&http.Cookie{Name:"id", Value:"123"})
		c.Response.String("")
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		}else {
			cookies := res.Cookies()
			if len(cookies) != 1 || cookies[0].Name != "id" || cookies[0].Value != "123" {
				t.Error(cookies)
			}
		}
	}()
}

func TestJSON(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	root := po.Root
	po.Root.Middleware(logRequest)
	user := testUser{
		Name:"吴浩麟",
		Age:23,
		Money:123.456,
		Alive:true,
		Notes:[]testNote{
			{Text:"明天去放风筝"},
			{Text:"今天我们去逛宜家啦"},
		},
	}
	root.Get("/hi", func(c *Context) {
		c.Response.JSON(user)
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		}else {
			res_bs, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			real_bs, _ := json.Marshal(user)
			if !reflect.DeepEqual(res_bs, real_bs) {
				t.Error(string(res_bs))
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != applicationJSONCharsetUTF8 {
				t.Error(ct)
			}
		}
	}()
}

func TestJSONP(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	root := po.Root
	po.Root.Middleware(logRequest)
	user := testUser{
		Name:"吴浩麟",
		Age:23,
		Money:123.456,
		Alive:true,
		Notes:[]testNote{
			{Text:"明天去放风筝"},
			{Text:"今天我们去逛宜家啦"},
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
		}else {
			res_bs, _ := ioutil.ReadAll(res.Body)
			res_str := string(res_bs)
			defer res.Body.Close()
			real_bs, _ := json.Marshal(user)
			real_str := string(real_bs)
			if res_str != ("hello(" + real_str + ")") {
				t.Error(string(res_str))
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != applicationJavaScriptCharsetUTF8 {
				t.Error(ct)
			}
		}
	}()
}

func TestXML(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	root := po.Root
	po.Root.Middleware(logRequest)
	user := testUser{
		Name:"吴浩麟",
		Age:23,
		Money:123.456,
		Alive:true,
		Notes:[]testNote{
			{Text:"明天去放风筝"},
			{Text:"今天我们去逛宜家啦"},
		},
	}
	root.Get("/hi", func(c *Context) {
		c.Response.XML(user)
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		}else {
			res_bs, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			real_bs, _ := xml.Marshal(user)
			if !reflect.DeepEqual(res_bs, real_bs) {
				t.Error(string(res_bs))
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != applicationXMLCharsetUTF8 {
				t.Error(ct)
			}
		}
	}()
}

func TestFile(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	root := po.Root
	po.Root.Middleware(logRequest)
	root.Get("/file", func(c *Context) {
		c.Response.File("test_resource/html/index.html")
	})
	defer func() {
		res, err := http.Get(baseURL + "/file")
		if err != nil {
			t.Error(err)
		}else {
			bs, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err == nil {
				str := string(bs)
				if str != "<h1>index.html</h1><b>{{.}}</b>" {
					t.Error(str)
				}
			}
		}
	}()
}

func TestString(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	root := po.Root
	po.Root.Middleware(logRequest)
	bodyStr := "hello,乓"
	root.Get("/hi", func(c *Context) {
		c.Response.String(bodyStr)
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		}else {
			bs, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
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
		}
	}()
}

func TestHTML(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	root := po.Root
	po.Root.Middleware(logRequest)
	bodyHTML := "<h1>hello,乓</>"
	root.Get("/hi", func(c *Context) {
		c.Response.HTML(bodyHTML)
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		}else {
			bs, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			html := string(bs)
			if html != bodyHTML {
				t.Error(html)
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != textHTMLCharsetUTF8 {
				t.Error(ct)
			}
		}
	}()
}

func TestRender(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	root := po.Root
	po.Root.Middleware(logRequest)
	po.LoadTemplateGlob("test_resource/html/*.html")
	root.Get("/render/:name", func(c *Context) {
		name := c.Request.Param("name")
		c.Response.Render(name, "中文")
	})
	defer func() {
		res, err := http.Get(baseURL + "/render/index.html")
		if err != nil {
			t.Error(err)
		}else {
			bs, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			html := string(bs)
			if html != "<h1>index.html</h1><b>中文</b>" {
				t.Error(html)
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != textHTMLCharsetUTF8 {
				t.Error(ct)
			}
		}
	}()
	defer func() {
		res, err := http.Get(baseURL + "/render/footer.html")
		if err != nil {
			t.Error(err)
		}else {
			bs, _ := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			html := string(bs)
			if html != "<h1>footer.html</h1><b>中文</b>" {
				t.Error(html)
			}
			ct := res.Header.Get(httpHeaderContentType)
			if ct != textHTMLCharsetUTF8 {
				t.Error(ct)
			}
		}
	}()
	go http.ListenAndServe(listenAddr, po)
}
