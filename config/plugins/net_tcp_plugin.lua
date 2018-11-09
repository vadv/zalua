package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
guage = require "prometheus_gauge"

local tcp_state_desc = {
  "established", "syn_sent", "syn_recv", "fin_wait1", "fin_wait2", "time_wait", "close",
  "close_wait", "last_ack", "listen", "closing"
}

local tcp_state_map = {}
for i,v in pairs(tcp_state_desc) do
  k = string.format("%02X", i)
  tcp_state_map[k] = v
end

-- регистрируем prometheus метрики
tcp_state = guage.new({
  help     = "system tcp state",
  namespace = "system",
  subsystem = "tcp",
  name      = "state",
  labels    = {"type"}
})

-- главный loop
while true do

  local result = {}; for k, _ in pairs(tcp_state_map) do result[k] = 0 end
  for line in io.lines("/proc/net/tcp") do
    local data = strings.split(line, " "); local state = data[5]
    if state and not(result[state] == nil) then result[state] = result[state] + 1 end
  end

  for k, v in pairs(result) do
    local t = tcp_state_map[k]
    tcp_state:set(v, {type=t})
  end

  time.sleep(60)
end
