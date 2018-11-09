package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
sysinfo = require "sysinfo"

-- регистрируем prometheus метрики
netstat = prometheus_gauge_labels.new({
  help     = "system tcp netstat",
  namespace = "system",
  subsystem = "tcp",
  name      = "netstat",
  labels    = { "type", "fqdn" }
})

-- главный loop
while true do

  local result, tcp_keys, tcp_vals = {}, {}, {}
  for line in io.lines("/proc/net/netstat") do
    if string.match(line, "TcpExt:") then
    if string.match(line, "TcpExt%:%s+%d") then tcp_vals = strings.split(line, " ") else tcp_keys = strings.split(line, " ") end
    end
  end

  for i, k in pairs(tcp_keys) do result[k] = tonumber(tcp_vals[i]) end
  for k, v in pairs(result) do
    if not(k == "TcpExt:") then
      local zabbix_key = "system.tcp.netstat["..k.."]"
      metrics.set_counter_speed(zabbix_key, v)
      local value = metrics.get(zabbix_key)
      if value then netstat:set({type=k, fqdn=sysinfo.fqdn}, tonumber(value)) end
    end
  end

  time.sleep(60)
end
