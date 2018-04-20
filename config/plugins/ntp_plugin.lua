if not(goruntime.goos == "linux") then return end

while true do

  local state, err = cmd.exec("ntpstat")
  if not(err == nil) then
    -- команда завершилась с ошибкой
    if not(state == nil) then
      -- есть stdout/stderr
      local msg = strings.split(state.stdout, "\n")[1]
      metrics.set("ntp.stat", "состояние синхронизации: "..tostring(msg))
    else
      -- state == nil означает что команда не найдена
      metrics.set("ntp.stat", "команда ntpstat не найдена")
      return
    end
  else
    -- команда выполнилась без ошибки
    metrics.set("ntp.stat", "ok")
  end

  time.sleep(60)
end
