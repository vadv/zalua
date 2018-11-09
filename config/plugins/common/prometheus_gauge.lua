-- загружаем дополнительные labels из соседнего файла additional_prometheus_labels.lua
-- если он присутствует
local current_dir = filepath.dir(debug.getinfo(1).source)
local additional_labels = {}
if goos.stat(current_dir.."/additional_prometheus_labels.lua") then
  additional_labels = require "additional_prometheus_labels"
end

-- http://lua-users.org/wiki/SimpleLuaClasses
local prom_gauge = {}
prom_gauge.__index = prom_gauge

function prom_gauge:new(options)
  local gauge = {}
  setmetatable(gauge, prom_gauge)
  gauge.options = options
  -- не забываем добавлять дополнительные labels
  for k, _ in pairs(additional_labels) do table.insert(options.labels, k) end
  gauge.additional_labels = additional_labels
  gauge.prometheus = prometheus_gauge_labels.new(options)
  return gauge
end

function prom_gauge:set(labels, value)
  -- немного удобства
  if not(type(labels) == "table") then
    labels, value = value, labels
  end
  -- не забываем добавлять дополнительные labels
  local real_labels = labels
  for k,v in pairs(self.additional_labels) do real_labels[k] = v end
  local prometheus = self.prometheus
  prometheus:set(real_labels, value)
end

return prom_gauge
