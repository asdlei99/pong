package pong

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// A Request represents an HTTP request received by a server or to be sent by a client.
// Request has some convenient method to get params form client
type Request struct {
	requestParamMap map[string]string
	//point to http.Request in golang's standard lib
	HTTPRequest *http.Request
}

// get Path param in request URL
//
// for example:
//	register router is
//		router.get("/user/:id",handle)
//	the request URL is
//		/user/123
//	request.Param("id") == "123"
// If key is not present, returns the empty string.
func (req *Request) Param(name string) string {
	return req.requestParamMap[name]
}

// get Query param in request URL
//
// for example:
//	the request URL is
//		/user?name=hal
//	request.Query("name") == "hal"
// If key is not present, returns the empty string.
func (req *Request) Query(name string) string {
	return req.HTTPRequest.URL.Query().Get(name)
}

// get Form param form request's body
//
// returns the first value for the named component of the POST
// or PUT request body. URL query parameters are ignored.
// support both application/x-www-form-urlencoded and multipart/form-data
//
// If key is not present, returns the empty string.
func (req *Request) Form(name string) string {
	return req.HTTPRequest.PostFormValue(name)
}

// returns the first file for the provided form key.
func (req *Request) File(name string) (multipart.File, *multipart.FileHeader, error) {
	return req.HTTPRequest.FormFile(name)
}

// parse request's body data as JSON and use standard lib json.Unmarshal to bind data to struct
//
// an error will return if json.Unmarshal return error
func (req *Request) BindJSON(pointer interface{}) error {
	bs, _ := ioutil.ReadAll(req.HTTPRequest.Body)
	return json.Unmarshal(bs, pointer)
}

// parse request's body data as XML and use standard lib XML.Unmarshal to bind data to struct
//
// an error will return if xml.Unmarshal return error
func (req *Request) BindXML(pointer interface{}) error {
	bs, _ := ioutil.ReadAll(req.HTTPRequest.Body)
	return xml.Unmarshal(bs, pointer)
}

// parse request's body post form as map and bind data to struct use filed name
//
// an error will return if the struct filed type is not support
func (req *Request) BindForm(pointer interface{}) error {
	ct := req.HTTPRequest.Header.Get(httpHeaderContentType)
	switch {
	case strings.HasPrefix(ct, applicationForm):
		req.HTTPRequest.ParseForm()
		return bind(pointer, req.HTTPRequest.Form)
	case strings.HasPrefix(ct, multipartForm):
		req.HTTPRequest.ParseMultipartForm(32 << 20) //32 MB
		return bind(pointer, req.HTTPRequest.Form)
	default:
		return ErrorTypeNotSupport
	}
}

// parse request's query params as map and bind data to struct use filed name
//
// an error will return if the struct filed type is not support
func (req *Request) BindQuery(pointer interface{}) error {
	m := req.HTTPRequest.URL.Query()
	return bind(pointer, m)
}

// auto bind will look request's http Header ContentType
//
// if request ContentType is applicationJSON will use BindJSON to parse
// if request ContentType is applicationXML will use BindXML to parse
// if request ContentType is applicationForm or multipartForm will use BindForm to parse
// else will return an ErrorTypeNotSupport error
func (req *Request) AutoBind(pointer interface{}) error {
	ct := req.HTTPRequest.Header.Get(httpHeaderContentType)
	switch {
	case strings.HasPrefix(ct, applicationJSON):
		return req.BindJSON(pointer)
	case strings.HasPrefix(ct, applicationXML):
		return req.BindXML(pointer)
	case strings.HasPrefix(ct, applicationForm):
		req.HTTPRequest.ParseForm()
		return bind(pointer, req.HTTPRequest.Form)
	case strings.HasPrefix(ct, multipartForm):
		req.HTTPRequest.ParseMultipartForm(32 << 20) //32 MB
		return bind(pointer, req.HTTPRequest.PostForm)
	default:
		return ErrorTypeNotSupport
	}
}

func bind(pointer interface{}, m map[string][]string) error {
	if pointer == nil {
		return errors.New("can't bind to nil")
	}
	typ := reflect.TypeOf(pointer)
	if typ.Kind() != reflect.Ptr {
		return errors.New("can only bind to pointer")
	}
	typ = typ.Elem()
	val := reflect.ValueOf(pointer).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}
		inputFieldName := typeField.Name
		inputValue, exists := m[inputFieldName]
		if !exists {
			continue
		}
		structFieldKind := structField.Kind()
		if numElems := len(inputValue); structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for i := 0; i < numElems; i++ {
				if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else {
			if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
				return err
			}
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return ErrorTypeNotSupport
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}
