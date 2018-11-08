local sysinfo = {}
sysinfo.fqdn = strings.trim(ioutil.readfile("/proc/sys/kernel/hostname"), "\n")

return sysinfo
