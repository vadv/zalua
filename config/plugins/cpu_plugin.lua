-- обработка строки из /proc/stat
function read_cpu_values(str)
  -- https://www.kernel.org/doc/Documentation/filesystems/proc.txt
  local fields = { "user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest", "guest_nice" }
  local row, offset = {}, 1
  for value in str:gmatch("(%d+)") do
    row[fields[offset]] = tonumber(value)
    offset = offset + 1
  end
  return row
end

-- главный loop
while true do
  for line in io.lines("/proc/stat") do

    -- разбираем строчку которая начинается с ^(cpu )
    local cpu_all_line = line:match("^cpu%s+(.*)")
    if cpu_all_line then
        local cpu_all_values = read_cpu_values(cpu_all_line)
        for key, value in pairs(cpu_all_values) do
          metrics.set_counter_speed("system.cpu.total."..key, value)
        end
    end

    -- вычисляем running, blocked
    local processes = line:match("^procs_(.*)")
    if processes then
      local key, val = string.match(processes, "^(%S+)%s+(%d+)")
      metrics.set("system.processes."..key, tonumber(val))
    end

    -- вычисляем context switching
    local ctxt = line:match("^ctxt (%d+)")
    if ctxt then metrics.set_counter_speed("system.cpu.ctxt", tonumber(ctxt)) end

  end
  utils.sleep(60)
end
