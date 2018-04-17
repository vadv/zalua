-- плагин который сообщает только одну метрику в zabbix
-- все процессы запущены и работают или какие-то не запущены или "флапают"
local services = {}

-- главный loop
while true do

  local problem = '' -- пытаемся сократить кол-во сообщений в zabbix

  -- если есть /etc/service
  if goos.stat('/etc/service') then
    for _, file in pairs(filepath.glob('/etc/service/*')) do

      local name = file:match('^/etc/service/(%S+)$')
      local run = (ioutil.readfile(file..'/supervise/stat') == "run\n")
      local uptime, stat = 0, goos.stat(file..'/supervise/pid')
      if stat then uptime = (time.unix() - stat.mod_time) end

      if run then
        if uptime < 60 then
          -- подозрительный сервис
          if services[name] and (services[name] < 60) then
            -- был до этого уже замечен, отмечаем как флапающий
            local desc = "'"..name.."' перезапускается"
            if problem == '' then problem = desc else problem = problem..', '..desc end
          end
        end
      else
        -- процесс слинкован, но не запущен и это уже плохо
        local desc = "'"..name.."' был слинкован в /etc/services/, но не запущен"
        if problem == '' then problem = desc else problem = problem..', '..desc end
      end
      services[name] = uptime

    end

    if problem == '' then problem = 'ok' else problem = 'Найдены проблемы с runit сервисами: '..problem end
    metrics.set('runit.problem', problem)
  else
    metrics.set('runit.problem', 'ok')
  end
  time.sleep(60)
end
