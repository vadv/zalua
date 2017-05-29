-- необходимы права 4777 на исполняемый go-файл

function read_proc_io()
  local row = {}
  for _,file in ipairs(filepath.glob("/proc/*/io")) do
    local data = ioutil.readfile(file)
    if data then
      for _, line in pairs(strings.split(data, "\n")) do
        local key, val = line:gmatch("^(%S+):%s+(%d+)$")
        if key then
          if not row[key] then row[key] = 0 end
          row[key] = row[key] + tonumber(val)
        end
      end
    end
  end
  return row
end

-- главный loop
while true do

  -- сбираем статистику за snapshot_timeout
  local snapshot_timeout = 10
  local row1 = read_proc_io()
  utils.sleep(snapshot_timeout)
  local row2 = read_proc_io()

  local logical_read = row2["rchar"] - row1["rchar"]
  local physical_read = row2["read_bytes"] - row1["read_bytes"]
  if (logical_read > 0) and (physical_read > 0) then
    metrics.set("system.io.read_hit", logical_read/physical_read)
  else
    metrics.set("system.io.read_hit", 0)
  end

  metrics.set_counter_speed("system.io.syscr", row1["syscr"])
  metrics.set_counter_speed("system.io.syscw", row1["syscw"])

  utils.sleep(60-snapshot_timeout)
end
