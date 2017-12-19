package dsl

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestDslYaml(t *testing.T) {

	testYamlLua := `
data = [[
a:
  b: 1
]]

t, err = yaml.decode(data)
if not(err == nil) then error(err) end
if not(t["a"]["b"] == 1) then error("not working, get: "..t["a"]["b"]) end
`

	state := lua.NewState()
	Register(NewConfig(), state)
	if err := state.DoString(testYamlLua); err != nil {
		t.Fatalf("execute lua error: %s\n", err.Error())
	}

}
