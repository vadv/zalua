local ignore_fs = "^(autofs|binfmt_misc|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs)$"
local ignore_mountpoints = "^/(dev|proc|sys|var/lib/docker/.+)($|/)"

-- главный loop
while true do
  local cpu_count = 0
  for line in io.lines("/proc/mounts") do
    local data = strings.split(line, " ")
    local device, mountpoint, fs = data[1], data[2], data[3]
    if not regexp.match(fs, ignore_fs) then
      if not regexp.match(mountpoint, ignore_mountpoints) then
        print("device: "..device.." mountpoint: "..mountpoint.." fs: "..fs)
        print("------------------")
      end
    end
  end
  time.sleep(10)
end
