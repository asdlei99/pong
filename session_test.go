package pong

import (
	"net/http"
	"testing"
	"net/http/cookiejar"
	"net/url"
	"reflect"
)

func TestNoSession(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	po.Root.Middleware(logRequest)
	po.EnableSession()
	root := po.Root
	root.Get("/hi", func(c *Context) {
		c.Response.String("")
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		}else {
			cookies := res.Cookies()
			if len(cookies) != 1 || cookies[0].Name != SessionCookiesName {
				t.Error(cookies)
			}
		}
	}()
}

func TestHasSession(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	po.Root.Middleware(logRequest)
	po.EnableSession()
	root := po.Root
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
		sUser := c.Session.Get("user")
		if !reflect.DeepEqual(sUser, user) {
			t.Error(sUser)
		}
	})
	defer func() {
		sid := po.SessionManager.NewSession()
		po.SessionManager.Write(sid, map[string]interface{}{"user":user})
		jar, _ := cookiejar.New(nil)
		url, _ := url.Parse(baseURL + "/hi")
		jar.SetCookies(url, []*http.Cookie{&http.Cookie{Name:SessionCookiesName, Value:sid}})
		client := http.Client{
			Jar:jar,
		}
		res, err := client.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		}else {
			cookies := res.Cookies()
			if len(cookies) != 1 || cookies[0].Name != SessionCookiesName {
				t.Error(cookies)
			}
		}
	}()
}