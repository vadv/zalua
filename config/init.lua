local current_file = debug.getinfo(1).source
local plugins_dir = filepath.dir(current_file).."/".."plugins"
local plugins = {}

-- (пе)загрузка конкретного плагина
function re_run_plugin_from_file(file)

  local metadata = plugins[file]
  local current_md5 = crypto.md5(ioutil.read_file(file))

  -- старт плагина
  if metadata == nil then
    log.info("run plugin: "..file)
    local p = plugin.new(file)
    metadata["plugin"] = p
    metadata["md5"] = current_md5
    p:run()
    plugins[file] = metadata
    return
  end

  -- если файл изменился - останавливаем старый и запускаем новый
  if not(metadata["md5"] == current_md5) then
    log.info("stop plugin: "..file.." with md5: "..metadata["md5"])
    metadata["plugin"]:stop()
    local p = plugin.new(file)
    metadata["plugin"] = p
    log.info("start plugin: "..file.." with md5: "..current_md5)
    p:run()
    plugins[file] = metadata
  end

end

-- запуск и остановка всех плагинов
function re_run_if_needed()

  local all_files = {}
  for file, _ in pairs(plugins) do all_files[file] = false end
  for _, file in pairs(filepath.glob(plugins_dir.."/*_plugin.lua")) do
    all_files[file] = true
    re_run_plugin_from_file(file)
  end

  -- нужно остановить те, что не найдены
  for file, found in pairs(plugins) do
    if not found then
      log.info("delete unknown plugin: "..file)
      local metadata = plugins[file]
      log.info("stop plugin: "..file.." with md5: "..metadata["md5"])
      metadata[file]:stop()
      table.remove(plugins, file)
    end
  end

end

-- супервизор для плагинов
while true do
  time.sleep(5)

  re_run_if_needed()
  local error_count = 0

  -- проверка статусов всех плагинов
  for file, metadata in pairs(plugins) do
    local p = metadata["plugin"]
    if not p:is_running() then
      local err = p:error()
      if err then
        -- плагин не запущен, и завершился с ошибкой
        log.error(err)
        log.info("start plugin: "..file.." with md5: "..metadata["md5"])
        p:run()
        error_count = error_count + 1
        metrics.set("zalua.error.last", err)
      else
        -- плагин остановлен и не завершился с ошибкой, удаляем его из списка
        table.remove(plugins, file)
      end
    end
  end
  metrics.set("zalua.error.count", error_count)

end
