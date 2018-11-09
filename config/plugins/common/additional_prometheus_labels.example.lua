local labels = {}
labels.fqdn = strings.trim(ioutil.readfile("/proc/sys/kernel/hostname"), "\n")

return labels
