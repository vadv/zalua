-- плагин который сообщает только одну метрику в zabbix
-- все процессы запущены и работают или какие-то не запущены или "флапают"
local services = {}

-- главный loop
while true do

  local bad_services = '' -- пытаемся сократить кол-во сообщений в zabbix

  -- если есть директория
  if os.stat('/etc/service') then
    for _, file in pairs(filepath.glob('/etc/service/*')) do

      local name = file:match('^/etc/service/(%S+)$')
      local run = (ioutil.readfile(file..'/supervise/stat') == "run\n")
      local uptime, stat = 0, os.stat(file..'/supervise/pid')
      if stat then uptime = (time.unix() - stat.mod_time) end

      if run then
        if uptime < 60 then
          -- подозрительный сервис
          if services[name] and (services[name] < 60) then
            -- был до этого уже замечен, отмечаем как флапающий
            local desc = "'"..name.."' has flapping uptime"
            if bad_services == '' then bad_services = desc else bad_services = bad_services..', '..desc end
          end
        end
      else
        -- процесс слинкован, но не запущен и это уже плохо
        local desc = "'"..name.."' has linked, but isn't running"
        if bad_services == '' then bad_services = desc else bad_services = bad_services..', '..desc end
      end
      services[name] = uptime

    end
  end

  if not (bad_services == '') then bad_services = 'Found problem with runit services: '..bad_services end
  metrics.set('runit.problem', bad_services)
  utils.sleep(60)
end
