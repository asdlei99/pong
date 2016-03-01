package pong

import (
	"testing"
	"net/http"
	"strings"
	"encoding/json"
	"bytes"
	"reflect"
	"encoding/xml"
	"net/url"
)

type testUser struct {
	Name  string
	Age   int
	Money float64
	Notes []testNote
	Alive bool
}

type testNote struct {
	Text string
}

func TestParam(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	root.Get("/user/:id/update/:data", func(c *Context) {
		if c.Request.Param("id") != "123" {
			t.Error()
		}
		if c.Request.Param("data") != "更新" {
			t.Error()
		}
		t.Log(`TestParam`)
	})
	defer http.Get(baseURL + "/user/123/update/更新")
	root.Post("/user/:id/update/:data", func(c *Context) {
		if c.Request.Param("id") != "123" {
			t.Error()
		}
		if c.Request.Param("data") != "更新" {
			t.Error()
		}
		t.Log(`TestParam`)
	})
	defer http.Post(baseURL + "/user/123/update/更新", textPlainCharsetUTF8, strings.NewReader(""))
}

func TestQuery(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	root.Get("/hi", func(c *Context) {
		id := c.Request.Query("id")
		if id != "吴浩麟" {
			t.Error(id)
		}
		t.Log(`TestQuery`)
	})
	defer http.Get(baseURL + "/hi?id=吴浩麟")
}

func TestForm(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	small := "吴浩麟"
	big := ""
	for i := 0; i < 10000; i++ {
		big += small
	}
	root.Post("/small", func(c *Context) {
		data := c.Request.Form("data")
		if data != small {
			t.Error(data)
		}
		t.Log(`TestForm`)
	})
	root.Post("/big", func(c *Context) {
		data := c.Request.Form("data")
		if data != big {
			t.Error(data)
		}
		t.Log(`TestForm`)
	})
	defer http.PostForm(baseURL + "/small", url.Values{
		"data":[]string{small},
	})
	defer http.PostForm(baseURL + "/big", url.Values{
		"data":[]string{big},
	})
}

func TestBindJSON(t *testing.T) {
	po, baseURL := runPong()
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
	root.Post("/hi", func(c *Context) {
		bindUser := testUser{}
		c.Request.BindJSON(&bindUser)
		if !reflect.DeepEqual(bindUser, user) {
			t.Error(bindUser, user)
		}
		t.Log(`TestBindJSON`)
	})
	bs, _ := json.Marshal(user)
	defer http.Post(baseURL + "/hi", applicationJSONCharsetUTF8, bytes.NewReader(bs))
}

func TestBindXML(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := testUser{
		Name:"吴浩麟",
		Age:23,
		Money:123.456,
		Alive:false,
		Notes:[]testNote{
			{Text:"明天去放风筝"},
			{Text:"今天我们去逛宜家啦"},
		},
	}
	root.Post("/hi", func(c *Context) {
		bindUser := testUser{}
		c.Request.BindXML(&bindUser)
		if !reflect.DeepEqual(bindUser, user) {
			t.Error(bindUser, user)
		}
		t.Log(`TestBindXML root.Post("/hi")`)
	})
	bs, _ := xml.Marshal(user)
	defer http.Post(baseURL + "/hi", applicationXMLCharsetUTF8, bytes.NewReader(bs))
}

func TestBindForm(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := testUser{
		Name:"吴浩麟",
		Age:23,
		Money:123.456,
		Alive:true,
	}
	root.Post("/hi", func(c *Context) {
		bindUser := testUser{}
		c.Request.BindForm(&bindUser)
		if !reflect.DeepEqual(bindUser, user) {
			t.Error(bindUser, user)
		}
		err := c.Request.BindForm(nil)
		if err == nil {
			t.Error("should return err")
		}
		t.Log(`TestBindForm root.Post("/hi")`)
	})
	defer http.PostForm(baseURL + "/hi", url.Values{
		"Name":[]string{"吴浩麟"},
		"Age":[]string{"23"},
		"Money":[]string{"123.456"},
		"Alive":[]string{"true"},
	})
}

func TestBindQuery(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := testUser{
		Name:"吴浩麟",
		Age:23,
		Money:123.456,
	}
	root.Get("/hi", func(c *Context) {
		bindUser := testUser{}
		c.Request.BindQuery(&bindUser)
		if !reflect.DeepEqual(bindUser, user) {
			t.Error(bindUser, user)
		}
		err := c.Request.BindQuery(nil)
		if err == nil {
			t.Error("should return err")
		}
		err = c.Request.BindQuery(user)
		if err == nil {
			t.Error("should return err")
		}
		t.Log(`TestBindQuery root.Get("/hi")`)
	})
	defer http.Get(baseURL + "/hi?Name=吴浩麟&Age=23&Money=123.456")
}

func TestAutoBind(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := testUser{
		Name:"吴浩麟",
		Age:23,
		Money:123.456,
	}
	root.Post("/hi", func(c *Context) {
		bindUser := testUser{}
		c.Request.AutoBind(&bindUser)
		if !reflect.DeepEqual(bindUser, user) {
			t.Error(bindUser, user)
		}
		err := c.Request.AutoBind(nil)
		if err == nil {
			t.Error("should return err")
		}
		err = c.Request.AutoBind(user)
		if err == nil {
			t.Error("should return err")
		}
		t.Log(`TestAutoBind`)
	})
	defer http.PostForm(baseURL + "/hi", url.Values{
		"Name":[]string{"吴浩麟"},
		"Age":[]string{"23"},
		"Money":[]string{"123.456"},
	})
	bs, _ := json.Marshal(user)
	defer http.Post(baseURL + "/hi", applicationJSONCharsetUTF8, bytes.NewReader(bs))
	bs, _ = xml.Marshal(user)
	defer http.Post(baseURL + "/hi", applicationXMLCharsetUTF8, bytes.NewReader(bs))
}