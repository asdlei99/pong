package pong

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"testing"
)

func TestNoSession(t *testing.T) {
	po, baseURL := runPong()
	po.EnableSession()
	root := po.Root
	root.Get("/hi", func(c *Context) {
		c.Response.String("")
	})
	defer func() {
		res, err := http.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		} else {
			cookies := res.Cookies()
			if len(cookies) != 1 || cookies[0].Name != SessionCookiesName {
				t.Error(cookies)
			}
			t.Log(`TestHasSession`)
		}
	}()
}

func TestHasSession(t *testing.T) {
	po, baseURL := runPong()
	po.EnableSession()
	root := po.Root
	user := testUser{
		Name:  "吴浩麟",
		Age:   23,
		Money: 123.456,
		Alive: true,
		Notes: []testNote{
			{Text: "明天去放风筝"},
			{Text: "今天我们去逛宜家啦"},
		},
	}
	root.Get("/hi", func(c *Context) {
		sUser := c.Session.Get("user")
		if !reflect.DeepEqual(sUser, user) {
			t.Error(sUser)
		}
		t.Log(`TestHasSession`)
	})
	defer func() {
		sid := po.SessionManager.NewSession()
		po.SessionManager.Write(sid, map[string]interface{}{"user": user})
		jar, _ := cookiejar.New(nil)
		url, _ := url.Parse(baseURL + "/hi")
		jar.SetCookies(url, []*http.Cookie{&http.Cookie{Name: SessionCookiesName, Value: sid}})
		client := http.Client{
			Jar: jar,
		}
		_, err := client.Get(url.String())
		if err != nil {
			t.Error(err)
		}
	}()
}

func TestUpdateSessionValue(t *testing.T) {
	po, baseURL := runPong()
	po.EnableSession()
	root := po.Root
	root.Get("/initSession", func(c *Context) {
		c.Session.Set("name", "吴浩麟")
		c.Session.Set("age", 23)
		c.Response.String("initSession")
	})
	root.Get("/updateSessionValue", func(c *Context) {
		c.Session.Set("name", "halwu")
		c.Session.Set("age", 100)
		c.Response.String("updateSessionValue")
	})
	defer func() {
		client := http.Client{}
		res, err := client.Get(baseURL + "/initSession")
		if err != nil {
			t.Error(err)
		} else {
			cookies := res.Cookies()
			if len(cookies) != 1 || cookies[0].Name != SessionCookiesName {
				t.Error(cookies)
			}
			sid := cookies[0].Value
			sValue := po.SessionManager.Read(sid)
			if sValue["name"] != "吴浩麟" || sValue["age"] != 23 {
				t.Error(sValue)
			}
			jar, _ := cookiejar.New(nil)
			url, _ := url.Parse(baseURL + "/updateSessionValue")
			jar.SetCookies(url, cookies)
			client.Jar = jar
			_, err = client.Get(url.String())
			if err != nil {
				t.Error(err)
			} else {
				sValue := po.SessionManager.Read(sid)
				if sValue["name"] != "halwu" || sValue["age"] != 100 {
					t.Error(sValue)
				}
			}
			t.Log(`TestUpdateSessionValue`)
		}
	}()
}

func TestResetSessionValue(t *testing.T) {
	po, baseURL := runPong()
	po.EnableSession()
	root := po.Root
	root.Get("/initSession", func(c *Context) {
		c.Session.Set("name", "吴浩麟")
		c.Session.Set("age", 23)
		c.Response.String("initSession")
	})
	root.Get("/resetSessionValue", func(c *Context) {
		c.Session.Reset()
		c.Response.String("resetSessionValue")
	})
	defer func() {
		client := http.Client{}
		res, err := client.Get(baseURL + "/initSession")
		if err != nil {
			t.Error(err)
		} else {
			cookies := res.Cookies()
			if len(cookies) != 1 || cookies[0].Name != SessionCookiesName {
				t.Error(cookies)
			}
			sid := cookies[0].Value
			jar, _ := cookiejar.New(nil)
			url, _ := url.Parse(baseURL + "/resetSessionValue")
			jar.SetCookies(url, cookies)
			client.Jar = jar
			res, err = client.Get(url.String())
			if err != nil {
				t.Error(err)
			} else {
				sid2 := res.Cookies()[0].Value
				if sid2 == sid {
					t.Error("sid should be diff")
				}
				if po.SessionManager.Read(sid) != nil {
					t.Error("old sid value in session store should nil")
				}
				sValue := po.SessionManager.Read(sid2)
				if sValue["name"] != "吴浩麟" || sValue["age"] != 23 {
					t.Error(sValue)
				}
			}
			t.Log(`TestResetSessionValue`)
		}
	}()
}

func TestDestorySession(t *testing.T) {
	po, baseURL := runPong()
	po.EnableSession()
	root := po.Root
	root.Get("/initSession", func(c *Context) {
		c.Session.Set("name", "吴浩麟")
		c.Session.Set("age", 23)
		c.Response.String("initSession")
	})
	root.Get("/destorySessionValue", func(c *Context) {
		c.Session.Destory()
		c.Response.String("destorySessionValue")
	})
	defer func() {
		client := http.Client{}
		res, err := client.Get(baseURL + "/initSession")
		if err != nil {
			t.Error(err)
		} else {
			cookies := res.Cookies()
			if len(cookies) != 1 || cookies[0].Name != SessionCookiesName {
				t.Error(cookies)
			}
			sid := cookies[0].Value
			jar, _ := cookiejar.New(nil)
			url, _ := url.Parse(baseURL + "/destorySessionValue")
			jar.SetCookies(url, cookies)
			client.Jar = jar
			res, err = client.Get(url.String())
			if err != nil {
				t.Error(err)
			} else {
				if po.SessionManager.Read(sid) != nil {
					t.Error("old sid value in session store should nil")
				}
				removeCookieHeader := res.Header.Get("Set-Cookie")
				if removeCookieHeader != SessionCookiesName+"=; Max-Age=0" {
					t.Error(removeCookieHeader)
				}
			}
			t.Log(`TestDestorySession`)
		}
	}()
}

func TestCheaterSession(t *testing.T) {
	po, baseURL := runPong()
	po.EnableSession()
	root := po.Root
	root.Get("/hi", func(c *Context) {
		c.Response.String("initSession")
	})
	defer func() {
		client := http.Client{}
		res, err := client.Get(baseURL + "/hi")
		if err != nil {
			t.Error(err)
		} else {
			jar, _ := cookiejar.New(nil)
			url, _ := url.Parse(baseURL + "/hi")
			cheaterSid := "cheaterSid-cheaterSid"
			jar.SetCookies(url, []*http.Cookie{&http.Cookie{Name: SessionCookiesName, Value: cheaterSid}})
			client.Jar = jar
			res, err = client.Get(url.String())
			if err != nil {
				t.Error(err)
			} else {
				cookies := res.Cookies()
				if len(cookies) != 1 || cookies[0].Name != SessionCookiesName {
					t.Error(cookies)
				}
			}
			t.Log(`TestCheaterSession`)
		}
	}()
}
