package pong

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
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
	defer http.Post(baseURL+"/user/123/update/更新", textPlainCharsetUTF8, strings.NewReader(""))
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
	defer http.PostForm(baseURL+"/small", url.Values{
		"data": []string{small},
	})
	defer http.PostForm(baseURL+"/big", url.Values{
		"data": []string{big},
	})
}

func TestBindJSON(t *testing.T) {
	po, baseURL := runPong()
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
	root.Post("/hi", func(c *Context) {
		bindUser := testUser{}
		c.Request.BindJSON(&bindUser)
		if !reflect.DeepEqual(bindUser, user) {
			t.Error(bindUser, user)
		}
		t.Log(`TestBindJSON`)
	})
	bs, _ := json.Marshal(user)
	defer http.Post(baseURL+"/hi", applicationJSONCharsetUTF8, bytes.NewReader(bs))
}

func TestBindXML(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := testUser{
		Name:  "吴浩麟",
		Age:   23,
		Money: 123.456,
		Alive: false,
		Notes: []testNote{
			{Text: "明天去放风筝"},
			{Text: "今天我们去逛宜家啦"},
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
	defer http.Post(baseURL+"/hi", applicationXMLCharsetUTF8, bytes.NewReader(bs))
}

func TestBindApplicationForm(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := testUser{
		Name:  "吴浩麟",
		Age:   23,
		Money: 123.456,
		Alive: true,
	}
	root.Post("/hi", func(c *Context) {
		bindUser := testUser{}
		c.Request.BindForm(&bindUser)
		if !reflect.DeepEqual(bindUser, user) {
			t.Error(bindUser, user)
		}
		c.Request.AutoBind(&bindUser)
		if !reflect.DeepEqual(bindUser, user) {
			t.Error(bindUser, user)
		}
		err := c.Request.BindForm(nil)
		if err == nil {
			t.Error("should return err")
		}
		t.Log(`TestBindForm root.Post("/hi")`)
	})
	defer http.PostForm(baseURL+"/hi", url.Values{
		"Name":  []string{"吴浩麟"},
		"Age":   []string{"23"},
		"Money": []string{"123.456"},
		"Alive": []string{"true"},
	})
}

func TestBindMultipartForm(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := testUser{
		Name:  "吴浩麟",
		Age:   23,
		Money: 123.456,
		Alive: true,
	}
	root.Post("/hi", func(c *Context) {
		bindUser := testUser{}
		c.Request.BindForm(&bindUser)
		if !reflect.DeepEqual(bindUser, user) {
			t.Error(bindUser, user)
		}
		c.Request.AutoBind(&bindUser)
		if !reflect.DeepEqual(bindUser, user) {
			t.Error(bindUser, user)
		}
		t.Log(`TestBindMultipartForm`)
	})
	defer func() {
		client := http.Client{}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("Name", "吴浩麟")
		writer.WriteField("Age", "23")
		writer.WriteField("Money", "123.456")
		writer.WriteField("Alive", "true")
		writer.Close()
		client.Post(baseURL+"/hi", writer.FormDataContentType(), body)
	}()
}

func TestBindQuery(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := testUser{
		Name:  "吴浩麟",
		Age:   23,
		Money: 123.456,
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
		Name:  "吴浩麟",
		Age:   23,
		Money: 123.456,
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
	defer http.PostForm(baseURL+"/hi", url.Values{
		"Name":  []string{"吴浩麟"},
		"Age":   []string{"23"},
		"Money": []string{"123.456"},
	})
	bs, _ := json.Marshal(user)
	defer http.Post(baseURL+"/hi", applicationJSONCharsetUTF8, bytes.NewReader(bs))
	bs, _ = xml.Marshal(user)
	defer http.Post(baseURL+"/hi", applicationXMLCharsetUTF8, bytes.NewReader(bs))
}

func TestUpdateFile(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	filePath := "test_resource/html/index.html"
	root.Post("/hi", func(c *Context) {
		file, header, err := c.Request.File("file")
		if err != nil {
			t.Error(err)
		}
		bs, err := ioutil.ReadAll(file)
		if err != nil {
			t.Error(err)
		}
		fileStr := string(bs)
		if header.Filename != filePath {
			t.Error(header.Filename)
		}
		if fileStr != "<h1>index.html</h1><b>{{.}}</b>" {
			t.Error(fileStr)
		}
		t.Log(`TestUpdateFile`)
	})
	defer func() {
		file, _ := os.Open(filePath)
		client := http.Client{}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		filePart, _ := writer.CreateFormFile("file", filePath)
		io.Copy(filePart, file)
		writer.Close()
		client.Post(baseURL+"/hi", writer.FormDataContentType(), body)
	}()
}

func TestBindContentTypeNotSupport(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	root.Post("/hi", func(c *Context) {
		user := &testUser{}
		err := c.Request.BindForm(user)
		if err != ErrorTypeNotSupport {
			t.Error(err)
		}
		err = c.Request.AutoBind(user)
		if err != ErrorTypeNotSupport {
			t.Error(err)
		}
		t.Log(`TestBindContentTypeNotSupport`)
	})
	defer func() {
		file, _ := os.Open("test_resource/html/index.html")
		_, err := http.Post(baseURL+"/hi", "ContentTypeNotSupport", file)
		if err != nil {
			t.Error(err)
		}
	}()
}

type fullType struct {
	Int        int
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Bool       bool
	Float32    float32
	Float64    float64
	String     string
	Slice      []string
	BoolSlice  []bool
	HandleFunc HandleFunc
}

func TestBind(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	full := fullType{
		Int:     123,
		Bool:    false,
		Slice:   []string{"1", "a", "中文"},
		Uint8:   12,
		Float32: 12.34,
	}
	root.Post("/hi", func(c *Context) {
		bindFull := fullType{}
		err := c.Request.BindForm(&bindFull)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(bindFull, full) {
			t.Error(bindFull, full)
		}
		t.Log(`TestBind`)
	})
	defer http.PostForm(baseURL+"/hi", url.Values{
		"Int":     []string{"123"},
		"Int8":    []string{""},
		"Int16":   []string{""},
		"Int32":   []string{""},
		"Int64":   []string{""},
		"Uint":    []string{""},
		"Uint8":   []string{"12"},
		"Uint16":  []string{""},
		"Uint32":  []string{""},
		"Uint64":  []string{""},
		"Bool":    []string{""},
		"Float32": []string{"12.34"},
		"Float64": []string{""},
		"String":  []string{""},
		"Slice":   []string{"1", "a", "中文"},
	})
	root.Post("/err", func(c *Context) {
		bindFull := fullType{}
		err := c.Request.BindForm(&bindFull)
		if err == nil {
			t.Error("should be err")
		}
		t.Log(`TestBind`)
	})
	defer http.PostForm(baseURL+"/err", url.Values{
		"Bool": []string{"abc"},
	})
}
