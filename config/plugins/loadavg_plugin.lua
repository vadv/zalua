package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
sysinfo = require "sysinfo"

-- регистрируем prometheus метрики
la1 = prometheus_gauge_labels.new({
  help     = "system loadavg 1m",
  namespace = "system",
  subsystem = "cpu",
  name      = "la_1m",
  labels    = { "fqdn" }
})

la5 = prometheus_gauge_labels.new({
  help     = "system loadavg 5m",
  namespace = "system",
  subsystem = "cpu",
  name      = "la_5m",
  labels    = { "fqdn" }
})

la15 = prometheus_gauge_labels.new({
  help     = "system loadavg 15m",
  namespace = "system",
  subsystem = "cpu",
  name      = "la_15m",
  labels    = { "fqdn" }
})

while true do
  local line = strings.trim(ioutil.readfile("/proc/loadavg"), "\n")
  local data = strings.split(line, " ")
  la1:set({fqdn=sysinfo.fqdn}, tonumber(data[1]))
  la5:set({fqdn=sysinfo.fqdn}, tonumber(data[2]))
  la15:set({fqdn=sysinfo.fqdn}, tonumber(data[3]))
  time.sleep(60)
end
