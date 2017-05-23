-- парсит содержимое /proc/diskstats
function diskstat()

  local result = {}
  -- https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
  local pattern = "(%d+)%s+(%d+)%s+(%S+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)%s+(%d+)"

  for line in io.lines("/proc/diskstats") do
    local major, minor, dev_name,
       rd_ios, rd_merges_or_rd_sec, rd_sec_or_wr_ios, rd_ticks_or_wr_sec,
       wr_ios, wr_merges, wr_sec, wr_ticks, ios_pgr, tot_ticks, rq_ticks = line:match(pattern)
    result[dev_name] = {
      major = tonumber(major), minor = tonumber(minor),
      rd_ios = tonumber(rd_ios), rd_merges_or_rd_sec = tonumber(rd_merges_or_rd_sec),
      rd_sec_or_wr_ios = tonumber(rd_sec_or_wr_ios), rd_ticks_or_wr_sec = tonumber(rd_ticks_or_wr_sec),
      wr_ios = tonumber(wr_ios), wr_merges = tonumber(wr_merges),
      wr_sec = tonumber(wr_sec), wr_ticks = tonumber(wr_ticks),
      ios_pgr = tonumber(ios_pgr), tot_ticks = tonumber(tot_ticks),
      rq_ticks = tonumber(rq_ticks)
    }
  end

  return result
end

-- главный loop
while true do

  local discovery = {}
  for dev_name, values in pairs(diskstat()) do
    -- собираем только sdX
    if string.match(dev_name, "sd[a-z]+$") then
      local discovery_item = {}; discovery_item["{#DEV}"] = dev_name; table.insert(discovery, discovery_item)
      metrics.set_speed("system.disk."..dev_name..".utilization", values.tot_ticks / 10)
      metrics.set_speed("system.disk."..dev_name..".read_bytes", values.rd_sec_or_wr_ios * 512) -- 2.6.32
      metrics.set_speed("system.disk."..dev_name..".read_ops", values.rd_ios)
      metrics.set_speed("system.disk."..dev_name..".write_bytes", values.wr_sec * 512)
      metrics.set_speed("system.disk."..dev_name..".write_ops", values.wr_ios)
    end
  end
  metrics.set("system.disk.discovery", json.encode({data = discovery}))

  utils.sleep(60)

end
