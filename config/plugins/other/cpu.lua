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

    -- разбираем строчку которая начинается с ^(cpu )
    local cpu_all_line = line:match("^cpu%s+(.*)")
    if cpu_all_line then
        local cpu_all_values = read_cpu_values(cpu_all_line)
        for key, value in pairs(cpu_all_values) do
          metrics.set_speed("system.cpu.all."..key, value)
        end
    end

    -- выясняем running, blocked `^procs_{running, blocked} \d+`
    local processes = line:match("^procs_(.*)")
    if processes then
      local key, val = string.match(processes, "^(%S+)%s+(%d+)")
      metrics.set("system.processes."..key, tonumber(val))
    end

    -- считаем cpu_count
    local number = line:match("^cpu(%d)%s+.*")
    if line:match("^cpu(%d)%s+.*") then
        cpu_number = tonumber(number) + 1
        if cpu_number > cpu_count then cpu_count = cpu_number end
    end

  end
  metrics.set("system.cpu.count", cpu_count)
  utils.sleep(60)
end
