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

-- /dev/sda => mountpoint, /dev/mapper/vg0-lv_slashroot => /
-- мы ищем только прямое совпадение!
function get_mountpoint_from_mounts(full_dev_name)
  for line in io.lines("/proc/mounts") do
    local reg_full_dev_name = full_dev_name:gsub("%-", "%S")
    local mountpoint = line:match("^"..reg_full_dev_name.."%s+(%S+)%s+")
    if mountpoint then return mountpoint end
  end
end

-- sdXD => mountpoint
function sd_mountpoint(sdX)
  -- нет задачи получить lvm, или md - мы это определяем ниже,
  -- тут мы ловим только напрямую подмонтированные
  return get_mountpoint_from_mounts("/dev/"..sdX)
end

-- dm-X => mountpoint
function dm_mountpoint(dmX)
  local name = ioutil.readfile("/sys/block/"..dmX.."/dm/name"):gsub("^%s+", ""):gsub("%s+$", "")
  if not name then return nil end
  return get_mountpoint_from_mounts("/dev/mapper/"..name)
end

-- mdX => mountpoint
function md_mountpoint(mdX)
  return get_mountpoint_from_mounts("/dev/"..mdX)
end

-- sd, md, dm => mountpoint
function get_mountpoint_by_dev(dev)
  if dev:match("^sd") then return sd_mountpoint(dev) end
  if dev:match("^dm") then return dm_mountpoint(dev) end
  if dev:match("^md") then return md_mountpoint(dev) end
end

-- mdX => raid0, raid1, ...
function md_level(mdX)
  return ioutil.readfile("/sys/block/"..mdX.."/md/level"):gsub("^%s+", ""):gsub("%s+$", "")
end

-- mdX => {sda = X, dm-0 = Y}
function md_device_sizes(mdX)
  local result = {}
  for _, path in pairs(filepath.glob("/sys/block/"..mdX.."/slaves/*")) do
    local dev = path:match("/sys/block/"..mdX.."/slaves/(%S+)$")
    result[dev] = tonumber(ioutil.readfile(path.."/size"))
  end
  return result
end

-- главный loop
while true do

  local devices_info, all_stats = {}, {}
  for dev, values in pairs(diskstat()) do
    if dev:match("^sd") or dev:match("^md") or dev:match("^dm") then
      local mountpoint = get_mountpoint_by_dev(dev)
      -- запоминаем только те, по которым мы нашли mountpoint
      if mountpoint then devices_info[dev] = {}; devices_info[dev]["mountpoint"] = mountpoint; end
      -- all stats мы заполняем для всех устройств, так как будет некое шульмование с mdX
      all_stats[dev] = {
        utilization = values.tot_ticks / 10, read_bytes = values.rd_sec_or_wr_ios * 512,
        read_ops = values.rd_ios, write_bytes = values.wr_sec * 512,
        write_ops = values.wr_ios
      }
    end
  end

  local discovery = {}
  -- теперь пришло время отослать собранные данные
  for dev, info in pairs(devices_info) do

    local mountpoint = info["mountpoint"]
    local discovery_item = {}; discovery_item["{#MOUNTPOINT}"] = mountpoint; table.insert(discovery, discovery_item)
    local utilization = 0

    if dev:match("sd") or dev:match("dm") then
      metrics.set_speed("system.disk.utilization["..mountpoint.."]", all_stats[dev]["utilization"])
    end

    -- а вот с md пошло шульмование про utilization
    if dev:match("md") then
      local slaves_info = md_device_sizes(dev)
      local total_slave_size = 0; for _, size in pairs(slaves_info) do total_slave_size = total_slave_size + size end
      local raid_level = md_level(dev)

      -- для raid{0,1} просчитываем utilization с весом
      -- вес высчитывается = (размер slave) / (сумму размера slave-устройств)
      if (raid_level == "raid0") or (raid_level == "raid1") then
        for slave, size in pairs(slaves_info) do
          local weight = size / total_slave_size
          utilization = utilization + (all_stats[slave]["utilization"] * weight)
        end
      end

      metrics.set_speed("system.disk.utilization["..mountpoint.."]", utilization)
    end

    -- остсылем все остальные метрики
    metrics.set_speed("system.disk.read_bytes_in_sec["..mountpoint.."]", all_stats[dev]["read_bytes"])
    metrics.set_speed("system.disk.read_ops_in_sec["..mountpoint.."]", all_stats[dev]["read_ops"])
    metrics.set_speed("system.disk.write_bytes_in_sec["..mountpoint.."]", all_stats[dev]["write_bytes"])
    metrics.set_speed("system.disk.write_ops_in_sec["..mountpoint.."]", all_stats[dev]["write_ops"])
  end

  metrics.set("system.disk.discovery", json.encode({data = discovery}))
  utils.sleep(60)
end
