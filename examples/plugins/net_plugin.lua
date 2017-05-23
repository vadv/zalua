proc_net_fields = {
  "rx.bytes", "rx.packets", "rx.errs", "rx.drop", "rx.fifo", "rx.frame", "rx.compressed", "rx.multicast",
  "tx.bytes", "tx.packets", "tx.errs", "tx.drop", "tx.fifo", "tx.frame", "tx.compressed", "tx.multicast",
}

-- обработка строки из /proc/net/dev без ethX:
function proc_net_field_values(str)
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
      for _, key in pairs(proc_net_fields) do
        local values = proc_net_field_values(row)
        metrics.set_speed("system.net."..interface.."."..key, values[key])
      end
    end
  end
  metrics.set("system.net.discovery", json.encode({data = discovery}))
  utils.sleep(60)
end
