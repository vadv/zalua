-- бесконечно слушаем указанный порт
prometheus.listen(":2345")

-- в случае падения нас перезапустят
while true do
  time.sleep(1)
end
