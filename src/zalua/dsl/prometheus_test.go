package dsl

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func TestPluginPrometheus(t *testing.T) {
	testLua := `
p = plugin.new("prometheus_test.lua")
p:run()
time.sleep(1)
err = p:error(); if err then error(err) end
time.sleep(1)
`
	state := lua.NewState()
	Register(NewConfig(), state)
	go func() {
		time.Sleep(time.Second)
		resp, err := http.Get("http://127.0.0.1:1111/metrics")
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		log.Printf("body:\n%s\n", body)
	}()
	if err := state.DoString(testLua); err != nil {
		t.Fatalf("execute lua error: %s\n", err.Error())
	}

}
