-- немного тупого "auto-discovery" для включения плагина
local enabled = false
if goos.stat('/var/lib/postgresql') then enabled = true end
if goos.stat('/var/lib/pgsql') then enabled = true end
if not enabled then return end

-- для работы данного плагина необходимо
-- 1. создать пользователя root `create user root with superuser;`
-- 2. дать разрешение в pg_hba: `local all root peer` не надо притворяться что unix-root не имеет полный доступ в базу
local connection = {
  host     = '/tmp',
  user     = 'root',
  database = 'postgres'
}
if goos.stat('/var/run/postgresql/.s.PGSQL.5432') then connection.host = '/var/run/postgresql' end

local previous_values = {}

-- открываем "главный" коннект
local main_db, err = postgres.open(connection)
if err then error(err) end

-- устанавливаем лимит на выполнение любого запроса 10s
local _, err = main_db:query("set statement_timeout to '10s'")
if err then main_db:close(); error(err) end

-- пробуем активировать расширение pg_stat_statements
local use_pg_stat_statements = false
local rows, err = main_db:query("select count(*) from pg_catalog.pg_extension where extname = 'pg_stat_statements'")
if not err then
  if rows[1][1] == 1 then
    use_pg_stat_statements = true
  else
    _, err = main_db:query("create extension pg_stat_statements;")
    if err then
      log.error("enable `pg_stat_statements`: "..tostring(err))
    else
      use_pg_stat_statements = true
    end
  end
end

while true do
  local discovery = {}

  -- in recovery (is slave)?
  local rows, err = main_db:query("select pg_catalog.pg_is_in_recovery() as pg_is_in_recovery")
  if err then main_db:close(); error(err) end
  local pg_is_in_recovery = rows[1][1]

  if pg_is_in_recovery then
    -- is slave
    local rows, err = main_db:query("select extract(epoch from pg_last_xact_replay_timestamp())")
    if err then main_db:close(); error(err) end
    -- брать во внимание что pg_last_xact_replay_timestamp это значение из времени мастера
    local replication_lag = time.unix() - rows[1][1] -- возможно clock_timestamp()
    metrics.set('postgres.wal.last_apply', replication_lag)
  else
    -- is master
    local _, err = main_db:query("select txid_current()")
    if err then main_db:close(); error(err) end
    metrics.set('postgres.wal.last_apply', 0)
    local rows, err = main_db:query("select pg_catalog.pg_xlog_location_diff (pg_catalog.pg_current_xlog_location(),'0/00000000')")
    if err then main_db:close(); error(err) end
    metrics.set_counter_speed('postgres.wal.write_bytes_in_sec', rows[1][1])
  end

  -- если активно pg_stat_statements
  if use_pg_stat_statements then
    local rows, err = main_db:query("select sum(total_time) as times, sum(calls) as calls, \
      sum(blk_read_time) as disk_read_time, sum(blk_write_time) as disk_write_time, sum(total_time - blk_read_time - blk_write_time) as other_time \
      from public.pg_stat_statements;")
    if not err then
      local current_times, current_calls, current_time = rows[1][1], rows[1][2], time.unix()
      metrics.set_counter_speed('postgres.time.disk_read_ms', rows[1][3])
      metrics.set_counter_speed('postgres.time.disk_write_ms', rows[1][4])
      metrics.set_counter_speed('postgres.time.other_ms', rows[1][5])
      local prev_times, prev_calls, prev_time = previous_values['total_time'], previous_values['calls'], previous_values['time']
      if prev_times then
        local diff_times, diff_calls, diff_time = (current_times - prev_times), (current_calls - prev_calls), (current_time - prev_time)
        if (diff_times > 0) and (diff_calls > 0) and (diff_time > 0) then
          metrics.set('postgres.queries.count', diff_calls/diff_time)
          metrics.set('postgres.queries.avg_time_ms', diff_times/diff_calls)
        end
      end
      previous_values['total_time'], previous_values['calls'], previous_values['time'] = current_times, current_calls, current_time
    end
  end

  -- < 9.6 только! (количество файлов pg_xlog)
  local rows, err = main_db:query("select count(*) from pg_catalog.pg_ls_dir('pg_xlog')")
  if not err then metrics.set('postgres.wal.count', rows[1][1]) end

  -- для >= 9.6 только
  local rows, err = main_db:query("select count(*) from pg_catalog.pg_stat_activity where wait_event is not null")
  if not err then metrics.set('postgres.connections.waiting', rows[1][1]) end

  -- кол-во autovacuum воркеров
  local rows, err = main_db:query("select count(*) from pg_catalog.pg_stat_activity where \
    query like '%autovacuum%' and state <> 'idle'")
  if not err then metrics.set('postgres.connections.autovacuum', rows[1][1]) end

  -- кол-во коннектов
  local rows, err = main_db:query("select state, count(*) from pg_catalog.pg_stat_activity group by state")
  if not err then
    for _, state in pairs({'active', 'idle', 'idle in transaction'}) do
      local state_value = 0
      -- если находим такой state в результатах, то присваеваем его
      for _, row in pairs(rows) do
        if (row[1] == state) then
          state_value = row[2]
        end
      end
      if state == 'idle in transaction' then state = 'idle_in_transaction' end
      metrics.set('postgres.connections.'..state, state_value)
    end
  end

  -- долго играющие запросы
  local long_queries_msg = "ok"
  local rows, err = main_db:query("select \
        s.query, extract(epoch from now() - s.query_start)::int as age \
    from pg_catalog.pg_stat_get_activity(null::integer) s \
      where not(s.state = 'idle') \
      and extract(epoch from now() - s.query_start)::int > 4500")
  if not err then
    local msg, count = "", 0
    for _, row in pairs(rows) do
      msg = msg .. " `" .. row[1] .. "` "
      count = count + 1
    end
    if count > 0 then long_queries_msg = "Найдено ".. tostring(count) .. " транзакци(й) которые запущены давно: ".. msg end
  end
  metrics.set('postgres.queries.long', long_queries_msg)

  -- кол-во локов
  local rows, err = main_db:query("select lower(mode), count(mode) FROM pg_catalog.pg_locks group by 1")
  if not err then
    for _, lock in pairs({'accessshare', 'rowshare', 'rowexclusive', 'shareupdateexclusive', 'share', 'sharerowexclusive', 'exclusive', 'accessexclusive'}) do
      local lock_full_name, lock_value = lock..'lock', 0
      for _, row in pairs(rows) do
        if (row[1] == lock_full_name) then
          lock_value = row[2]
        end
      end
      metrics.set('postgres.locks.'..lock, lock_value)
    end
  end

  -- хит по блокам и статистика по строкам
  local rows, err = main_db:query("select sum(blks_hit) as blks_hit, sum(blks_read) as blks_read, \
      sum(tup_returned) as tup_returned, sum(tup_fetched) as tup_fetched, sum(tup_inserted) as tup_inserted, \
      sum(tup_updated) as tup_updated, sum(tup_deleted) as tup_deleted \
    from pg_catalog.pg_stat_database")
  if not err then
    local current_hit, current_read = rows[1][1], rows[1][2]
    local prev_hit, prev_read = previous_values['blks_hit'], previous_values['blks_read']
    if prev_hit then
      local hit_rate, diff_hit, diff_read = 100, (current_hit - prev_hit), (current_read - prev_read)
      if (diff_hit > 0) and (diff_read > 0) then
        hit_rate = 100*diff_hit/(diff_read+diff_hit)
      end
      metrics.set('postgres.blks.hit_rate', hit_rate)
    end
    metrics.set_counter_speed('postgres.blks.hit', current_hit)
    metrics.set_counter_speed('postgres.blks.read', current_read)
    previous_values['blks_hit'], previous_values['blks_read'] = current_hit, current_read
    metrics.set_counter_speed('postgres.rows.returned', rows[1][3])
    metrics.set_counter_speed('postgres.rows.fetched', rows[1][4])
    metrics.set_counter_speed('postgres.rows.inserted', rows[1][5])
    metrics.set_counter_speed('postgres.rows.updated', rows[1][6])
    metrics.set_counter_speed('postgres.rows.deleted', rows[1][7])
  end

  -- выполняем из главной базы общий запрос на размеры и статусы
  local rows, err = main_db:query("select \
      datname, pg_catalog.pg_database_size(datname::text), pg_catalog.age(datfrozenxid) \
      from pg_catalog.pg_database where datistemplate = false")
  if err then main_db:close(); error(err) end
  for _, row in pairs(rows) do
    local dbname, size, age = row[1], row[2], row[3]
    metrics.set('postgres.database.size['..dbname..']', size)
    metrics.set('postgres.database.age['..dbname..']', age)
    local discovery_item = {}; discovery_item["{#DATABASE}"] = dbname; table.insert(discovery, discovery_item)
  end
  metrics.set("postgres.database.discovery", json.encode({data = discovery}))
  time.sleep(60)
end
