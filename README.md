# ZaLua

## Цели:

* легко расширяемый мониторинг написаный на скриптовом языке
* использовать транспорт zabbix-agent

## Решение:

толстый go-бинарь с встроенным lua-интерпретатором

## Запуск и получение метрик

При запуске проверяется что по сокету `/tmp/zalua-mon.sock` отвечает сервер на ping-pong сообщение,
иначе запускается демон-сервер который начинает слушать этот сокет.
демон пытает запустить `/etc/zalua/config.lua` dsl которого описан ниже.

### Доступные команды клиента:

```
-v, -version, --version
        Get version
-k, -kill, --kill, kill
        Kill server
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
$ zalua -p
/etc/zalua/plugins/buddyinfo_plugin.lua         false           /etc/zalua/plugins/buddyinfo_plugin.lua at 1: parse error
/etc/zalua/plugins/numa_plugin.lua              false           <no error>
/etc/zalua/plugins/snmp_plugin.lua              true            <no error>


$ zalua -m
system.tcp.active               0.000000                1495702422
system.tcp.established          1.000000                1495702422
system.tcp.failed               0.000000                1495702422
system.tcp.passive              0.000000                1495702422
system.tcp.resets               0.000000                1495702422
system.disk.read_bytes[/video]          23454385.653432         1495748781
system.disk.read_ops[/video]            46.463452               1495748781
system.disk.utilization[/video]         51.029993               1495748781
system.disk.write_bytes[/video]         6872556.411429          1495748781
system.disk.write_ops[/video]           25.724433               1495748781

$ zalua -g system.tcp.active
0.000000
```

### Пример UserParameter:

```
UserParameter=disk.utilization[*], /usr/bin/zalua -g system.disk.$1.utilization
```

## lua-DSL

* *plugin*:
    * `p = plugin.new(filename)` создать плагин
    * `p:run()` запустить плагин
    * `p:stop()` остановить вызвав ошибку stop в плагине
    * `p:is_running()` запущен или нет плагин
    * `p:error()` текст последний ошибки или nil

* *metric*:
    * `metric.set(key, val, <ttl>)` установить значение метрики key, val может быть string, number. ttl по дефолту 300 секунд
    * `metric.set_speed(key, val, <ttl>)` тоже самое, но считает скорость измерения
    * `metric.set_counter_speed(key, val, <ttl>)` тоже самое, но считает только положительную скорость измерения
    * `metric.get(key)` получить значение метрики key
    * `metric.list()` список метрик
    * `metric.delete(key)` удалить значение метрики key

* *utils*:
    * `utils.sleep(N)` проспать N секунд

* *ioutil*:
    * `ioutil.readfile(filename)` вернуть содержимое файла

* *strings*:
    * `strings.split(str, delim)` порт golang strings.split()

* *json*:
    * `json.encode(N)` lua-table в string
    * `json.decode(N)` string в lua-table

* *filepath*:
    * `filepath.base(filename)` порт golang filepath.Base()
    * `filepath.dir(filename)` порт golang filepath.Dir()
    * `filepath.ext(filename)` порт golang filepath.Ext()
    * `filepath.glob(mask)` порт golang filepath.Glob(), в случае ошибки возращает nil.

* *os*:
    * `os.stat(filename)` os.stat возвращает таблицу с полями `{size, is_dir, mod_time}`, в случае ошибки возращает nil.
    * `os.pagesize()` возвращет pagesize

* *log*:
    * `log.error(msg)` сообщение в лог с уровнем error
    * `log.info(msg)` с уровнем info

## Примеры плагинов

### Diskstat

![await](/img/await.png)

* пытается сопоставить блочному девайсу /mount/pount
* расчитывает await и utilization по тем блочным девайсам, по которым ядро не ведет статистику (mdraid)
