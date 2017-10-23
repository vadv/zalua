if not goos.stat('/sys/devices/system/node/node1/numastat') then return end

-- основной loop
while true do

  local numastats = {total={}}

  local paths = filepath.glob('/sys/devices/system/node/node*/numastat')
  for _, path in pairs(paths) do
    local node_num = path:match('/node(%d+)/')
    numastats[node_num] = {}
    for line in io.lines(path) do
      local key, val = line:match('(%S+)%s+(%d+)')
      numastats[node_num][key] = tonumber(val)
      if numastats.total[key] == nil then numastats.total[key] = 0 end
      numastats.total[key] = numastats.total[key] + val
    end
  end

  local discovery = {}
  for num, stats in pairs(numastats) do
    local discovery_item = {}; discovery_item["{#NUMA_NODE}"] = num; table.insert(discovery, discovery_item)
    metrics.set_counter_speed('system.numa.numa_hit['..num..']', stats['numa_hit']) -- A process wanted to allocate memory from this node and succeeded.
    metrics.set_counter_speed('system.numa.numa_miss['..num..']', stats['numa_miss']) -- A process wanted to allocate memory from another node but ended up with memory from this node
    metrics.set_counter_speed('system.numa.numa_foreign['..num..']', stats['numa_foreign']) -- A process wanted to allocate memory from this node but ended up with memory from another node
    metrics.set_counter_speed('system.numa.interleave_hit['..num..']', stats['interleave_hit']) -- Interleaving wanted to allocate memory from this node and succeeded
    metrics.set_counter_speed('system.numa.local_node['..num..']', stats['local_node']) -- A process ran on this node and got memory from it
    metrics.set_counter_speed('system.numa.other_node['..num..']', stats['other_node']) -- A process ran on this node and got memory from another node
  end
  metrics.set("system.numa.discovery", json.encode({data = discovery}))

  time.sleep(60)
end
