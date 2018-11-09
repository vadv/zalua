package.path = filepath.dir(debug.getinfo(1).source)..'/common/?.lua;'.. package.path
sysinfo = require "sysinfo"

local ignore_fs = "^(autofs|binfmt_misc|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs)$"
local ignore_mountpoints = "^/(dev|proc|sys|var/lib/docker/.+)($|/)"

-- регистрируем prometheus метрики
bytes_free = prometheus_gauge_vec.new({
  help     = "system filesystem percent free bytes",
  namespace = "system",
  subsystem = "disk",
  name      = "percent_bytes_free",
  vec       = { "device", "mountpoint", "fs", "options", "fqdn" }
})

inodes_free = prometheus_gauge_vec.new({
  help     = "system filesystem percent free inodes",
  namespace = "system",
  subsystem = "disk",
  name      = "percent_inodes_free",
  vec       = { "device", "mountpoint", "fs", "options", "fqdn" }
})

total_expose = prometheus_gauge_vec.new({
  help     = "system filesystem expose",
  namespace = "system",
  subsystem = "disk",
  name      = "expose",
  vec       = { "type", "device", "mountpoint", "fs", "options", "fqdn" }
})

-- главный loop
while true do
  local cpu_count = 0
  for line in io.lines("/proc/mounts") do

    local data = strings.split(line, " ")
    local device, mountpoint, fs, options = data[1], data[2], data[3], data[4]

    if not regexp.match(fs, ignore_fs) then
      if not regexp.match(mountpoint, ignore_mountpoints) then
        local info, err = syscall.statfs(mountpoint)
        if err then
          log.error("get syscall.statfs for "..mountpoint.." error: "..err)
        else
          if not (info["size"] == 0) then
            local free = 100 - (info["free"]/info["size"])*100
            bytes_free:set({device=device, mountpoint=mountpoint, fs=fs, options=options, fqdn=sysinfo.fqdn}, free)
          end
          if not (info["files"] == 0 ) then
            local free = 100 - (info["files_free"]/info["files"])*100
            inodes_free:set({device=device, mountpoint=mountpoint, fs=fs, options=options, fqdn=sysinfo.fqdn}, free)
          end
          for k, v in pairs(info) do
            total_expose:set({type=k, device=device, mountpoint=mountpoint, fs=fs, options=options, fqdn=sysinfo.fqdn}, v)
          end
        end
      end
    end

  end
  time.sleep(60)
end
