local enabled = false
if goos.stat('/etc/service/pgbouncer-6432') then enabled = true end
if not enabled then return end

local connection = {
  host = '127.0.0.1',
  port = 6432,
  user = 'postgres',
  database = 'pgbouncer',
}
local pgbouncer_db, err = postgres.open(connection)
if err then error(err) end

while true do

  local rows, err = pgbouncer_db:query("show stats")
  if err then pgbouncer_db:close(); error(err) end
  local total_query_count, total_received, total_sent, total_query_time, total_wait_time = 0, 0, 0, 0, 0
  for _, row in pairs(rows) do
    total_query_count = total_query_count + tonumber(row[3])
    total_received = total_received + tonumber(row[4])
    total_sent = total_sent + tonumber(row[5])
    total_query_time = total_query_time + tonumber(row[7])
    total_wait_time = total_wait_time + tonumber(row[8])
  end
  metrics.set_counter_speed('pgbouncer.query.count', total_query_count)
  metrics.set_counter_speed('pgbouncer.query.total_time', total_query_time/1000000)
  metrics.set_counter_speed('pgbouncer.query.wait_time', total_wait_time/1000000)
  metrics.set_counter_speed('pgbouncer.sent', total_sent)
  metrics.set_counter_speed('pgbouncer.received', total_received)

  time.sleep(60)
end
