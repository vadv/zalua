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
  local data = ioutil.readfile("/sys/block/"..mdX.."/md/level")
  if data then
    return data:gsub("%s+$", "")
  else
    return nil
  end
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

-- вычисляем значения которые зависят от предыдущих значений
calc_values = {}
function calc_value(dev, values)
  if calc_values[dev] == nil then calc_values[dev] = {} end
  if calc_values[dev]["data"] == nil then calc_values[dev]["data"] = {} end
  -- проставляем при первом проходе previous
  if calc_values[dev]["data"]["previous"] == nil then calc_values[dev]["data"]["previous"] = values; return; end

  local previous, current = calc_values[dev]["data"]["previous"], values

  -- вычисляем await https://github.com/sysstat/sysstat/blob/v11.5.6/common.c#L816
  local ticks = ((current.rd_ticks_or_wr_sec - previous.rd_ticks_or_wr_sec) + (current.wr_ticks - previous.wr_ticks))
  local io_sec = (current.rd_ios + current.wr_ios) - (previous.rd_ios + previous.wr_ios)
  if (io_sec > 0) and (ticks > 0) then
    calc_values[dev]["await"] = ticks / io_sec
  else
    -- игнорируем проворот счетчика
    if (io_sec == 0) then calc_values[dev]["await"] = 0 end
  end

  -- перетираем предыдущее значение
  calc_values[dev]["data"]["previous"] = values
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
        utilization = values.tot_ticks / 10,
        read_bytes = values.rd_sec_or_wr_ios * 512, read_ops = values.rd_ios,
        write_bytes = values.wr_sec * 512, write_ops = values.wr_ios
      }
      calc_value(dev, values)
    end
  end

  local discovery = {}
  -- теперь пришло время отослать собранные данные
  for dev, info in pairs(devices_info) do

    local mountpoint = info["mountpoint"]
    local discovery_item = {}; discovery_item["{#MOUNTPOINT}"] = mountpoint; table.insert(discovery, discovery_item)
    local utilization, await = 0, nil

    if dev:match("^sd") or dev:match("^dm") then
      utilization = all_stats[dev]["utilization"]
      await = calc_values[dev]["await"]
    end

    -- а вот с md пошло шульмование про utilization
    if dev:match("^md") then
      local slaves_info = md_device_sizes(dev)
      local total_slave_size = 0; for _, size in pairs(slaves_info) do total_slave_size = total_slave_size + size end
      local raid_level = md_level(dev)
      if raid_level then -- пропускаем непонятный raid
        -- для raid{0,1} просчитываем utilization с весом
        -- вес высчитывается = (размер slave) / (сумму размера slave-устройств)
        if (raid_level == "raid0") or (raid_level == "raid1") then
          for slave, size in pairs(slaves_info) do
            local weight = size / total_slave_size
            utilization = utilization + (all_stats[slave]["utilization"] * weight)
            local slave_await = calc_values[slave]["await"]
            if slave_await then
              if await == nil then await = 0 end
              await = await + (slave_await * weight)
            end
          end
        end
      end
    end

    metrics.set_counter_speed("system.disk.utilization["..mountpoint.."]", utilization)
    if await then metrics.set("system.disk.await["..mountpoint.."]", await) end

    -- остсылем все остальные метрики
    metrics.set_counter_speed("system.disk.read_bytes_in_sec["..mountpoint.."]", all_stats[dev]["read_bytes"])
    metrics.set_counter_speed("system.disk.read_ops_in_sec["..mountpoint.."]", all_stats[dev]["read_ops"])
    metrics.set_counter_speed("system.disk.write_bytes_in_sec["..mountpoint.."]", all_stats[dev]["write_bytes"])
    metrics.set_counter_speed("system.disk.write_ops_in_sec["..mountpoint.."]", all_stats[dev]["write_ops"])
    metrics.set_counter_speed("system.disk.all_ops_in_sec["..mountpoint.."]", all_stats[dev]["read_ops"] + all_stats[dev]["write_ops"])
  end

  metrics.set("system.disk.discovery", json.encode({data = discovery}))
  time.sleep(60)
end
