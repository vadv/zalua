package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
guage = require "prometheus_gauge"

-- регистрируем prometheus метрики
la1 = guage:new({
  help     = "system loadavg 1m",
  namespace = "system",
  subsystem = "cpu",
  name      = "la_1m"
})

la5 = guage:new({
  help     = "system loadavg 5m",
  namespace = "system",
  subsystem = "cpu",
  name      = "la_5m"
})

la15 = guage:new({
  help     = "system loadavg 15m",
  namespace = "system",
  subsystem = "cpu",
  name      = "la_15m"
})

while true do
  local line = strings.trim(ioutil.readfile("/proc/loadavg"), "\n")
  local data = strings.split(line, " ")
  la1:set(tonumber(data[1]))
  la5:set(tonumber(data[2]))
  la15:set(tonumber(data[3]))
  time.sleep(60)
end
