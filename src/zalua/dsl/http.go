package dsl

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const USER_AGENT = `zalua`

func (d *dslConfig) dslHttpGet(L *lua.LState) int {
	url := L.CheckString(1)
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("http create request: %s\n", err.Error())))
		return 2
	}
	req.Header.Set("User-Agent", USER_AGENT)
	response, err := client.Do(req)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("http error: %s\n", err.Error())))
		return 2
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("http read response error: %s\n", err.Error())))
		return 2
	}
	// write response
	result := L.NewTable()
	L.SetField(result, "code", lua.LNumber(response.StatusCode))
	L.SetField(result, "body", lua.LString(string(data)))
	L.SetField(result, "url", lua.LString(url))
	L.Push(result)
	return 1
}

func (d *dslConfig) dslHttpEscape(L *lua.LState) int {
	query := L.CheckString(1)
	escapedUrl := url.QueryEscape(query)
	L.Push(lua.LString(escapedUrl))
	return 1
}

func (d *dslConfig) dslHttpUnEscape(L *lua.LState) int {
	query := L.CheckString(1)
	url, err := url.QueryUnescape(query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("unescape error: %s\n", err.Error())))
		return 2
	}
	L.Push(lua.LString(url))
	return 1
}
