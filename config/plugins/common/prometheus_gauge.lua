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
  options.labels = options.labels or {}
  -- не забываем добавлять дополнительные labels
  for k, _ in pairs(additional_labels) do table.insert(options.labels, k) end
  gauge.additional_labels = additional_labels
  gauge.prometheus = prometheus_gauge_labels.new(options)
  return gauge
end

function prom_gauge:set(value, labels)
  -- немного удобства
  if not(type(value) == "number") then labels, value = value, labels end
  labels = labels or {}
  -- не забываем добавлять дополнительные labels
  local real_labels = labels
  for k,v in pairs(self.additional_labels) do real_labels[k] = v end
  self.prometheus:set(real_labels, value)
end

function prom_gauge:set_from_metrics(labels, key)
  -- немного удобства
  if not(type(labels) == "table") then labels, key = key, labels end
  local value = metrics.get(key) -- nil or string
  if value then self:set(labels, tonumber(value)) end
end

return prom_gauge
