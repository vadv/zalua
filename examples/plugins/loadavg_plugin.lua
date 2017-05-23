function loadavg()
  local pattern = "(%d*.%d+)%s+(%d*.%d+)%s+(%d*.%d+)%s+(%d+)/(%d+)%s+(%d+)"

  local file = io.open("/proc/loadavg")
  local content = file:read("*a")
  file:close()

  local minute1avg, minute5avg, minute15avg, runnable, exist, lastpid = string.match(content, pattern)
  return { minute1avg = tonumber(minute1avg), minute5avg = tonumber(minute5avg),
    minute15avg = tonumber(minute15avg), runnable = tonumber(runnable),
    exist = tonumber(exist), lastpid = tonumber(lastpid) }
end

while true do
  local avg = loadavg()
  metrics.set("system.processes.total", avg.exist)
--  metrics.set("system.la.1", avg.minute1avg)
--  metrics.set("system.la.5", avg.minute5avg)
--  metrics.set("system.la.15", avg.minute15avg)
  metrics.set_counter_speed("system.fork_rate", avg.lastpid)
  utils.sleep(60)
end
