package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
guage = require "prometheus_gauge"

directions = {"rx", "tx"}
types = {"bytes", "packets", "errs", "drop", "fifo", "frame", "compressed", "multicast"}
proc_net_fields = {}
for _, direction in pairs(directions) do
  for _, t in pairs(types) do
    table.insert(proc_net_fields, direction.."."..t)
  end
end

-- регистрируем prometheus метрики
gauge_net = guage:new({
  help     = "system net info",
  namespace = "system",
  subsystem = "net",
  name      = "info",
  labels    = { "direction", "interface", "type" }
})

-- обработка строки из /proc/net/dev без ethX:
function proc_net_field_value(str)
  local row, offset = {}, 1
  for value in str:gmatch("(%d+)") do
    row[proc_net_fields[offset]] = tonumber(value)
    offset = offset + 1
  end
  return row
end

-- основной loop
while true do
  local discovery = {}
  for line in io.lines("/proc/net/dev") do
    local interface, row = line:match("(%S+):%s+(.+)$")
    if not (interface == nil) and (interface:match("^vlan") or interface:match("^bond") or interface:match("^eth")) then
      local discovery_item = {}; discovery_item["{#DEV}"] = interface; table.insert(discovery, discovery_item)
      for _, direction in pairs(directions) do
        for _, t in pairs(types) do
          local key = direction.."."..t
          local value = proc_net_field_value(row)
          metrics.set_speed("system.net."..key.."["..interface.."]", value[key])
          gauge_net:set_from_metrics("system.net."..key.."["..interface.."]", {direction = direction, type = t, interface = interface})
        end
      end
    end
  end
  metrics.set("system.net.discovery", json.encode({data = discovery}))
  time.sleep(60)
end
