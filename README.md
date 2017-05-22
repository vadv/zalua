# ZaLua

## Цели:

* легко расширяемый мониторинг написаный на скриптовом языке
* использовать транспорт zabbix-agent

## Решение:

толстый go-бинарь с встроенным lua-интерпретатором

## Запуск и получение метрик

При запуске проверяется что по сокету /tmp/zalua-mon.sock отвечает сервер на ping-pong сообщение,
иначе запускается демон-сервер который начинает слушать этот сокет.
демон пытает запустить /etc/zalua/config.lua dsl которого описан ниже.

### Доступные команды:

```
-v, -version, --version
        Get version
-m, -metrics, --list-metrics, metrics
        List of known metrics
-p, -plugins, --plugins, plugins
        List of running plugins
-g, -get, --get, --get-metric, get <metric>
        Get metric value
-ping, --ping, ping
        Ping pong game
```

### Пример вывода:

```
$ zalua -m
system.disk.discovery           {"data":[{"{#DEV}":"sda"},{"{#DEV}":"sdb"},{"{#DEV}":"sdc"}]}           1495559213
system.disk.sda.read_bytes              11582.010802            1495559213
system.disk.sda.read_ops                0.784406                1495559213
system.disk.sda.utilization             0.470755                1495559213
system.disk.sda.write_bytes             123693.144210           1495559213
system.disk.sda.write_ops               3.078738                1495559213
system.disk.sdb.read_bytes              1481291.405964          1495559213
system.disk.sdb.read_ops                23.665508               1495559213
system.disk.sdb.utilization             2.568763                1495559213
system.disk.sdb.write_bytes             1412481.968397          1495559213
system.disk.sdb.write_ops               5.769718                1495559213
system.disk.sdc.read_bytes              1474332.771937          1495559213
system.disk.sdc.read_ops                23.214414               1495559213
system.disk.sdc.utilization             2.486211                1495559213
system.disk.sdc.write_bytes             1411871.836263          1495559213
system.disk.sdc.write_ops               5.561949                1495559213
```

## lua-DSL

* *plugin*:
    * `p = plugin.new(filename)` создать плагин
    * `p.run()` запустить
    * `p.check()` вернет ошибку если плагин не запущен или завершился с ошибкой
    * `p.stop()` остановить

* *metric*:
    * `metric.set(key, val, <ttl>)` установить значение метрики key, val может быть string, number. ttl по дефолту 300 секунд
    * `metric.set_speed(key, val, <ttl>)` тоже самое, но считает скорость измерения
    * `metric.get(key)` получить значение метрики key
    * `metric.list()` список метрик
    * `metric.delete(key)` удалить значение метрики key

* *utils*:
    * `utils.sleep(N)` проспать N секунд

* *json*:
    * `json.encode(N)` lua-table в string
    * `json.decode(N)` string в lua-table

* *filepath*:
    * `filepath.base(filename)` порт golang filepath.Base()
    * `filepath.dir(filename)` порт golang filepath.Dir()
    * `filepath.ext(filename)` порт golang filepath.Ext()
    * `filepath.glob(mask)` порт golang filepath.Glob(), может вызвать ошибку.

* *log*:
    * `log.error(msg)` сообщение в лог с уровнем error
    * `log.info(msg)` с уровнем info
