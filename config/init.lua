local current_file = debug.getinfo(1).source
local plugins_dir = filepath.dir(current_file).."/".."plugins"
local plugins = {} -- {filename= {plugin = p, md5 = md5}, ...}
local try_plugins = {} -- {filename = count_of_try}

-- удаление плагина
function delete_plugin(file)
  local metadata = plugins[file]
  if not(metadata == nil) then metadata["plugin"]:stop() end
  -- пересоздаем плагины, проще способа удалить по ключу не нашел
  local new_plugins = {}
  for old_file, old_metadata in pairs(plugins) do
    if not(old_file == file) then new_plugins[old_file] = old_metadata end
  end
  plugins = new_plugins
end

-- (пере)загрузка конкретного плагина
function re_run_plugin_from_file(file)

  local metadata = plugins[file]
  local current_md5 = crypto.md5(ioutil.readfile(file))

  -- старт плагина
  if metadata == nil then
    metadata = {}
    local p = plugin.new(file)
    metadata["plugin"] = p
    metadata["md5"] = current_md5
    p:run()
    plugins[file] = metadata
    return
  end

  -- если файл изменился - останавливаем старый и запускаем новый
  if not(metadata["md5"] == current_md5) then
    metadata["plugin"]:stop()
    local p = plugin.new(file)
    metadata["plugin"] = p
    metadata["md5"] = current_md5
    p:run()
    plugins[file] = metadata
  end

end

-- запуск и остановка плагинов
function re_run_if_needed()

  local found_files = {}
  for file, _ in pairs(plugins) do found_files[file] = false end

  for _, file in pairs(filepath.glob(plugins_dir.."/*_plugin.lua")) do
    found_files[file] = true
    re_run_plugin_from_file(file)
  end

  -- нужно остановить те, что не найдены
  for file, found in pairs(found_files) do
    if not found then
      local metadata = plugins[file]
      metadata["plugin"]:stop()
      delete_plugin(file)
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
        p:run()
        error_count = error_count + 1
        metrics.set("zalua.error.last", err)
      else
        -- плагин остановлен и не завершился с ошибкой
        -- попробуем его запустить позднее через перезапуск (после 10 попыток)
        local try_count = try_plugins[file]
        if try_count == nil then try_count = 0 end
        try_count = try_count + 1
        -- отправляем на рестарт
        if try_count > 10 then
          delete_plugin(file)
          try_plugins[file] = 0
        end

      end
    end
  end
  metrics.set("zalua.error.count", error_count)

end
