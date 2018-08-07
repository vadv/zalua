if not goos.stat('/var/log/messages') then return end

local syslog_layout = "Jan  2 15:04:05 2006"

while true do

  local min_time = os.time()-(24*60*60)
  local count_oom, count_segfault = 0, 0

  local scanner = tac.open("/var/log/messages")
  while true do
    local line = scanner:line()
    if line == nil then break end
    -- отрезаем первые символы в которых находиться время и прибавляем год
    local time_value = line:sub(0, 15)
    time_value = time_value .. " ".. os.date("%Y")
    -- пытаемся распарсить
    local log_time, err = time.parse(syslog_layout, time_value)
    -- выходим, если распаршенное время меньше min_time
    if err == nil then if log_time < min_time then break end end
    if string.find(line, "Out of memory: Kill process %d+ (%S+)") then count_oom = count_oom + 1 end
    if string.find(line, "kernel: (%S+)%[%d+%]: segfault at ") then count_segfault = count_segfault + 1 end
  end
  scanner:close()

  local messages = "ok"
  if count_oom + count_segfault > 0 then
    messages = "Найдено проблем с OOM: " .. count_oom .. ", проблем с segfault: ".. count_segfault .." за последние 24 часов."
  end
  metrics.set("system.messages.problem", messages, 10*60)


  time.sleep(5 * 60)
end
