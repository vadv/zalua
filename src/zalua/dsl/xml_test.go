package dsl

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestDslXml(t *testing.T) {

	testXmlLua := `

data = [[
<channels>
    <channel id="1" xz1="600" />
    <channel id="2"           />
    <channel id="x" xz2="600" />
</channels>
]]

t, err = xmlpath.parse(data, "/channels/channel/@id")
if not(err == nil) then error(err) end
x = 0
for k,v in pairs(t) do
  x = x + 1
  print("node "..v)
end

if not(x == 3) then error("x: "..x) end

`

	state := lua.NewState()
	Register(NewConfig(), state)
	if err := state.DoString(testXmlLua); err != nil {
		t.Fatalf("execute lua error: %s\n", err.Error())
	}

}
