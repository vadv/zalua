package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
guage = require "prometheus_gauge"

local ignore_fs = "^(autofs|binfmt_misc|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs)$"
local ignore_mountpoints = "^/(dev|proc|sys|var/lib/docker/.+)($|/)"

-- регистрируем prometheus метрики
bytes_free = guage:new({
  help     = "system filesystem percent free bytes",
  namespace = "system",
  subsystem = "disk",
  name      = "percent_bytes_free",
  labels    = { "device", "mountpoint", "fs" }
})

inodes_free = guage:new({
  help     = "system filesystem percent free inodes",
  namespace = "system",
  subsystem = "disk",
  name      = "percent_inodes_free",
  labels    = { "device", "mountpoint", "fs" }
})

total_expose = guage:new({
  help     = "system filesystem expose",
  namespace = "system",
  subsystem = "disk",
  name      = "expose",
  labels    = { "type", "device", "mountpoint", "fs" }
})

-- главный loop
while true do
  for line in io.lines("/proc/mounts") do

    local data = strings.split(line, " ")
    local device, mountpoint, fs = data[1], data[2], data[3]

    if not regexp.match(fs, ignore_fs) then
      if not regexp.match(mountpoint, ignore_mountpoints) then
        local info, err = syscall.statfs(mountpoint)
        if err then
          log.error("get syscall.statfs for "..mountpoint.." error: "..err)
        else
          if not (info["size"] == 0) then
            local free = 100 - (info["free"]/info["size"])*100
            bytes_free:set(free, {device=device, mountpoint=mountpoint, fs=fs})
          end
          if not (info["files"] == 0 ) then
            local free = 100 - (info["files_free"]/info["files"])*100
            inodes_free:set(free, {device=device, mountpoint=mountpoint, fs=fs})
          end
          for k, v in pairs(info) do
            total_expose:set(v, {type=k, device=device, mountpoint=mountpoint, fs=fs})
          end
        end
      end
    end

  end
  time.sleep(60)
end
