local current_file = debug.getinfo(1).source
local plugins_dir = filepath.dir(current_file).."/".."plugins"
local plugins = {}

-- загружаем все плагины
for _, file in pairs(filepath.glob(plugins_dir.."/*_plugin.lua")) do
  local p = plugin.new(file)
  table.insert(plugins, p)
  p:run()
end

-- супервизор для плагинов
while true do
  utils.sleep(5)

  for num, p in pairs(plugins) do
    if not p:is_running() then
      local err = p:error()
      if err then
        -- плагин не запущен, и завершился с ошибкой
        log.error(err)
        p:run()
      else
        -- плагин остановлен и не завершился с ошибкой, удаляем его из списка
        table.remove(plugins, num)
      end
    end
  end

end
