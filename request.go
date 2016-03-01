package pong

import (
	"net/http"
	"strings"
	"encoding/json"
	"encoding/xml"
	"reflect"
	"errors"
	"strconv"
	"io/ioutil"
	"mime/multipart"
)

type Request struct {
	requestParamMap map[string]string
	HTTPRequest     *http.Request
}

func (req *Request)Param(name string) string {
	return req.requestParamMap[name]
}

func (req *Request)Query(name string) string {
	return req.HTTPRequest.URL.Query().Get(name)
}

func (req *Request)Form(name string) string {
	return req.HTTPRequest.PostFormValue(name)
}

func (req *Request)File(name string) (multipart.File, *multipart.FileHeader, error) {
	return req.HTTPRequest.FormFile(name)
}

func (req *Request)BindJSON(pointer interface{}) error {
	bs, _ := ioutil.ReadAll(req.HTTPRequest.Body)
	return json.Unmarshal(bs, pointer)
}

func (req *Request)BindXML(pointer interface{}) error {
	bs, _ := ioutil.ReadAll(req.HTTPRequest.Body)
	return xml.Unmarshal(bs, pointer)
}

func (req *Request)BindForm(pointer interface{}) error {
	ct := req.HTTPRequest.Header.Get(httpHeaderContentType)
	switch {
	case strings.HasPrefix(ct, applicationForm):
		req.HTTPRequest.ParseForm()
		return bind(pointer, req.HTTPRequest.Form)
	case strings.HasPrefix(ct, multipartForm):
		req.HTTPRequest.ParseMultipartForm(32 << 20)//32 MB
		return bind(pointer, req.HTTPRequest.Form)
	default:
		return ErrorTypeNotSupport
	}
}

func (req *Request)BindQuery(pointer interface{}) error {
	m := req.HTTPRequest.URL.Query()
	return bind(pointer, m)
}

func (req *Request)AutoBind(pointer interface{}) error {
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
		req.HTTPRequest.ParseMultipartForm(32 << 20)//32 MB
		return bind(pointer, req.HTTPRequest.Form)
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