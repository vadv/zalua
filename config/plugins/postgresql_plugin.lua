-- немного тупого "auto-discovery" для включения плагина
local enabled = false
if os.stat('/var/lib/postgresql') then enabled = true end
if os.stat('/var/lib/pgsql') then enabled = true end
if not enabled then return end

-- для работы данного плагина необходимо
-- 1. создать пользователя root `create user root;`
-- 2. дать разрешение в pg_hba: `local all root peer`
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
