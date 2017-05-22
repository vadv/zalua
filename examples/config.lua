local current_file = debug.getinfo(1).source
local plugins_dir = filepath.dir(current_file).."/".."plugins"
local plugins = {}

-- загружаем все плагины
for _, file in pairs(filepath.glob(plugins_dir.."/*_plugin.lua")) do
    local p = plugin.new(file)
    p:run()
    table.insert(plugins, p)
end

-- супервизор для плагинов
while true do
    utils.sleep(10)
    for _, p in pairs(plugins) do
        local running, err = pcall(function() p:check() end)
        if not running then log.error(err); p:run() end
    end
end
