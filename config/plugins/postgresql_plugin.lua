-- немного тупого "auto-discovery" для включения плагина
local enabled = false
if os.stat('/var/lib/postgresql') then enabled = true end
if os.stat('/var/lib/pgsql') then enabled = true end
if not enabled then return end

-- для работы данного плагина необходимо
-- 1. создать пользователя root `create user root with superuser;`
-- 2. дать разрешение в pg_hba: `local all root peer` не надо притворяться что unix-root не имеет полный доступ в базу
local connection = {
  host     = '/tmp',
  user     = 'root',
  database = 'postgres'
}
if os.stat('/var/run/postgresql/.s.PGSQL.5432') then connection.host = '/var/run/postgresql' end

-- открываем "главный" коннект
local main_db, err = postgres.open(connection)
if err then error(err) end

-- устанавливаем лимит на выполнение любого запроса 10s
local _, err = main_db:query("set statement_timeout to '10s'")
if err then error(err) end

while true do
  local discovery = {}

  -- in recovery (is slave)?
  local rows, err = main_db:query("select pg_catalog.pg_is_in_recovery() as pg_is_in_recovery")
  if err then error(err) end
  local pg_is_in_recovery = rows[1][1]

  if not pg_is_in_recovery then
    -- is master
    local rows, err = main_db:query("select pg_catalog.pg_xlog_location_diff \
      (pg_catalog.pg_current_xlog_location(),'0/00000000')")
    if err then error(err) end
    metrics.set_counter_speed('postgres.wal.write', rows[1][1])
  end

  -- < 9.6 only!
  local rows, err = main_db:query("select count(*) from pg_catalog.pg_ls_dir('pg_xlog')")
  if not err then metrics.set('postgres.wal.count', rows[1][1]) end

  -- выполняем из главной базы общий запрос на размеры и статусы
  local rows, err = main_db:query("select \
      datname, pg_catalog.pg_database_size(datname::text), pg_catalog.age(datfrozenxid) \
      from pg_catalog.pg_database where datistemplate = false")
  if err then error(err) end
  for _, row in pairs(rows) do
    local dbname, size, age = row[1], row[2], row[3]
    metrics.set('postgres.database.size['..dbname..']', size)
    metrics.set('postgres.database.age['..dbname..']', age)
    local discovery_item = {}; discovery_item["{#DATABASE}"] = dbname; table.insert(discovery, discovery_item)
  end
  metrics.set("postgres.database.discovery", json.encode({data = discovery}))
  time.sleep(60)
end
