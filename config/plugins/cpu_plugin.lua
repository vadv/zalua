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
  local cpu_count = 0
  for line in io.lines("/proc/stat") do

    -- считаем cpu_count
    local number = line:match("^cpu(%d+)%s+.*")
    if number then
      number = tonumber(number) + 1; if number > cpu_count then cpu_count = number end
    end

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
    if ctxt then metrics.set_counter_speed("system.stat.ctxt", tonumber(ctxt)) end

    -- вычисляем processes
    local processes = line:match("^processes (%d+)")
    if processes then metrics.set_counter_speed("system.processes.fork_rate", tonumber(processes)) end

    -- вычисляем interupts
    local intr = line:match("^intr (%d+)")
    if intr then metrics.set_counter_speed("system.stat.intr", tonumber(intr)) end

  end
  metrics.set("system.cpu.count", cpu_count)
  time.sleep(60)
end
