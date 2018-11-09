package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
guage = require "prometheus_gauge"

-- регистрируем prometheus метрики
netstat = guage.new({
  help="system tcp netstat", namespace="system", subsystem="tcp",
  name="netstat", labels={"type"}
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
      netstat:set(zabbix_key, {type=k})
    end
  end

  time.sleep(60)
end
