package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
guage = require "prometheus_gauge"

function read_ss()
  local state, err = cmd.exec("/usr/sbin/ss --summary")
  if err == nil then
    -- команда не завершилась с ошибкой
    if not(state == nil) then
      -- есть stdout/stderr
      local result = {}
      for _, line in pairs(strings.split(state.stdout, "\n")) do
        if line:match("^TCP: ") then
          result["established"] = tonumber(line:match("estab (%d+),"))
          result["closed"] = tonumber(line:match("closed (%d+),"))
          result["orphaned"] = tonumber(line:match("orphaned (%d+),"))
          result["synrecv"] = tonumber(line:match("synrecv (%d+),"))
          local tw1, tw2 = line:match("timewait (%d+)/(%d+)")
          result["timewait"] = tonumber(tw1) + tonumber(tw2)
          return result
        end
      end
    end
  end
end

tcp_state = guage:new({ name = "system_tcp_state", labels = {"type"} })

-- главный loop
while true do
  local result = read_ss()
  if result then for k, v in pairs(result) do tcp_state:set({type=k}, v) end end
  time.sleep(60)
end
