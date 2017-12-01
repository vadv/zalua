function process(file)
  local result = {}
  for line in io.lines(file) do
    -- MemTotal:        8060536 kB
    local key, value = line:match("(%S+)%:%s+%d+%s+kB"), line:match("%S%:%s+(%d+)%s+kB")
    if (not(key == nil)) and (not(value == nil)) then result[key] = tonumber(value*1024) end
  end
  return result
end

-- основной loop
while true do
  local row = process("/proc/meminfo")
  local total, free, cached = 0, 0, 0
  for key, val in pairs(row) do
    if key == "MemFree" then
      free = val
      metrics.set("sys.memory.free", tostring(val))
    elseif key == "MemTotal" then
      total = val
    elseif key == "MemShared" then
      metrics.set("sys.memory.shared", tostring(val))
    elseif key == "Buffers" then
      metrics.set("sys.memory.buffers", tostring(val))
    elseif key == "Cached" then
      cached = val
      metrics.set("sys.memory.cached", tostring(val))
    end
  end
  metrics.set("sys.memory.other", tostring(total-free-cached))
  time.sleep(60)
end
