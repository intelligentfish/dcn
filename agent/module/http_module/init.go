package http_module

import (
	"bytes"
	"encoding/json"
	lua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	factoryMap map[string]func(L *lua.LState) *lua.LFunction
)

// HttpGet http get method
func HttpGet(token, urlStr string) (code int, raw []byte, err error) {
	var urlObj *url.URL
	if urlObj, err = url.Parse(urlStr); nil != err {
		return
	}
	// TODO URL black list && white list
	var res *http.Response
	if res, err = http.Get(urlStr); nil != err {
		return
	}
	if nil != res {
		code = res.StatusCode
		if nil != res.Body {
			defer res.Body.Close()
			raw, err = ioutil.ReadAll(res.Body)
		}
	}
	_ = urlObj
	//fmt.Println("token", token)
	//fmt.Println("rawQuery", urlObj.RawQuery)
	//fmt.Println("raw", raw)
	return
}

// HttpPost http post method
func HttpPost(token, urlStr, contentType string,
	req map[string]string) (code int, raw []byte, err error) {
	var urlObj *url.URL
	if urlObj, err = url.Parse(urlStr); nil != err {
		return
	}
	// TODO URL black list && white list
	const contentTypeJson = "application/json"
	var res *http.Response
	switch strings.ToLower(contentType) {
	case contentTypeJson:
		if raw, err = json.Marshal(req); nil != err {
			raw = nil
			return
		}
		res, err = http.Post(urlStr, contentTypeJson, bytes.NewReader(raw))
	default:
		data := url.Values{}
		for k, v := range req {
			data.Set(k, v)
		}
		res, err = http.PostForm(urlStr, data)
	}
	if nil != res {
		code = res.StatusCode
		if nil != res.Body {
			defer res.Body.Close()
			raw, err = ioutil.ReadAll(res.Body)
		}
	}
	_ = urlObj
	//fmt.Println("token", token)
	//fmt.Println("rawQuery", urlObj.RawQuery)
	//fmt.Println("raw", raw)
	return
}

// setLuaHTTPResult set lua http method execute result
func setLuaHTTPResult(L *lua.LState, code int, raw []byte, err error) int {
	L.Push(lua.LNumber(code))
	L.Push(lua.LString(string(raw)))
	if nil == err {
		L.Push(lua.LString(""))
	} else {
		L.Push(lua.LString(err.Error()))
	}
	return 3
}

// httpGetFactory http get method factory
func httpGetFactory(L *lua.LState) *lua.LFunction {
	return L.NewFunction(func(L *lua.LState) int {
		token := L.ToString(1)
		urlStr := L.ToString(2)

		code, raw, err := HttpGet(token, urlStr)

		return setLuaHTTPResult(L, code, raw, err)
	})
}

// httpPostFactory http post method factory
func httpPostFactory(L *lua.LState) *lua.LFunction {
	return L.NewFunction(func(L *lua.LState) int {
		token := L.ToString(1)
		urlStr := L.ToString(2)
		contentType := L.ToString(3)
		args := L.ToTable(4)
		req := make(map[string]string, 0)
		args.ForEach(func(key lua.LValue, value lua.LValue) {
			req[key.String()] = value.String()
		})

		code, raw, err := HttpPost(token, urlStr, contentType, req)

		return setLuaHTTPResult(L, code, raw, err)
	})
}

// Inject inject object into lua state
func Inject(L *lua.LState) {
	for k, v := range factoryMap {
		L.SetGlobal(k, v(L))
	}
}

// init initialize method
func init() {
	// factory method map
	factoryMap = map[string]func(L *lua.LState) *lua.LFunction{
		"httpGet":  httpGetFactory,
		"httpPost": httpPostFactory,
	}
}
