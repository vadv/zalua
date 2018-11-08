function process(file)
  local result = {}
  for line in io.lines(file) do
    -- MemTotal:        8060536 kB
    local key, value = line:match("(%S+)%:%s+%d+%s+kB"), line:match("%S%:%s+(%d+)%s+kB")
    if (not(key == nil)) and (not(value == nil)) then result[key] = tonumber(value*1024) end
  end
  return result
end

-- регистрируем prometheus метрики
gauge_memory = prometheus_gauge_vec.new({
  help     = "system memory discovery",
  namespace = "system",
  subsystem = "memory",
  name      = "bytes",
  vec       = { "type" }
})


-- основной loop
while true do
  local row = process("/proc/meminfo")
  local total, free, cached, shared, buffers = 0, 0, 0, 0, 0
  for key, val in pairs(row) do
    if key == "MemFree" then
      free = val
    elseif key == "MemTotal" then
      total = val
    elseif key == "MemShared" then
      shared = val
    elseif key == "Buffers" then
      buffers = val
    elseif key == "Cached" then
      cached = val
    end
    gauge_memory:set({ type = key }, val)
  end
  metrics.set("system.memory.free", tostring(free))
  metrics.set("system.memory.cached", tostring(cached))
  metrics.set("system.memory.shared", tostring(shared))
  metrics.set("system.memory.buffers", tostring(buffers))
  metrics.set("system.memory.other", tostring(total-free-cached))
  time.sleep(60)
end
