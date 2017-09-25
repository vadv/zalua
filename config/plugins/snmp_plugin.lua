-- парсит tcp строку из /proc/net/snmp
function parse_tcp(str)

  local tcp_fields = {
    'RtoAlgorithm', 'RtoMin', 'RtoMax',
    'MaxConn', 'ActiveOpens', 'PassiveOpens',
    'AttemptFails', 'EstabResets', 'CurrEstab',
    'InSegs', 'OutSegs', 'RetransSegs', 'InErrs',
    'OutRsts', 'InCsumErrors'}

  local row, offset = {}, 1
  for value in str:gmatch("(%d+)") do
    row[tcp_fields[offset]] = tonumber(value)
    offset = offset + 1
  end

  return row

end


-- основной loop
while true do

  for line in io.lines('/proc/net/snmp') do
    -- парсим tcp статистику
    local tcp_data_line = line:match('Tcp:%s+(%d+.+)$')
    if tcp_data_line then
      local row = parse_tcp(tcp_data_line)
      metrics.set_counter_speed('system.tcp.active', row['ActiveOpens']) -- The number of active TCP openings per second
      metrics.set_counter_speed('system.tcp.passive', row['PassiveOpens']) -- The number of passive TCP openings per second
      metrics.set_counter_speed('system.tcp.failed', row['AttemptFails']) -- The number of failed TCP connection attempts per second
      metrics.set_counter_speed('system.tcp.resets', row['EstabResets']) -- The number of TCP connection resets
      metrics.set_counter_speed('system.tcp.retransmit', row['RetransSegs'])
      metrics.set('system.tcp.established', row['CurrEstab']) -- The number of currently open connections
    end
  end

  time.sleep(60)
end
