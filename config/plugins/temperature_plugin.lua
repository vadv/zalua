local enabled = false
for _, file in ipairs(filepath.glob("/sys/bus/platform/devices/coretemp*/temp*_input")) do
  log.info("found coretemp file: "..file..", start plugin")
  enabled = true
  break
end
if not enabled then return end

function get_max_temp()
  local data = {}
  for _, file in ipairs(filepath.glob("/sys/bus/platform/devices/coretemp*/temp*_input")) do
    local temp = tonumber(ioutil.readfile(file))/1000
    table.insert(data, temp)
  end
  table.sort(data)
  return data[#data]
end

-- главный loop
while true do
  local maxt = get_max_temp()
  metrics.set("system.cpu.temp", maxt)

  time.sleep(10)
end
