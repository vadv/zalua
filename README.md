# ZaLua

## Цели:

* легко расширяемый мониторинг написаный на скриптовом языке
* использовать транспорт zabbix-agent

## Решение:

толстый go-бинарь с встроенным lua-интерпретатором

## Запуск и получение метрик

При запуске проверяется что по сокету `/tmp/zalua-mon.sock` отвечает сервер на ping-pong сообщение,
иначе запускается демон-сервер который начинает слушать этот сокет.
демон пытает запустить `/etc/zalua/init.lua` dsl которого описан ниже.

### Доступные команды клиента:

```
-v, -version, --version
        Get version
-e, --execute-file, execute file (without server)
        Execute dsl from file (for testing case)
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
/etc/zalua/plugins/bad_plugin.lua    false           /etc/zalua/plugins/buddyinfo_plugin.lua at 1: parse error
/etc/zalua/plugins/numa_plugin.lua   false           <no error>
/etc/zalua/plugins/snmp_plugin.lua   true            <no error>


$ zalua -m
system.disk.read_bytes[/video]          23454385.65         1495748781
system.disk.read_ops[/video]            46.46               1495748781
system.disk.utilization[/video]         51.02               1495748781
system.disk.write_bytes[/video]         6872556.41          1495748781
system.disk.write_ops[/video]           25.72               1495748781
runit.problem           Found problem with runit services: 'nginx' has linked, but isn't running   1495748781

$ zalua -g system.tcp.active
0.00
```

### Пример UserParameter:

```
UserParameter=disk.utilization[*], /usr/bin/zalua -g system.disk.$1.utilization
```

## lua-DSL

* *plugin*:
    * `p = plugin.new(filename)` загрузить плагин
    * `p:run()` запустить плагин
    * `p:stop()` остановить вызвав ошибку stop в плагине
    * `p:is_running()` запущен или нет плагин
    * `p:error()` текст последний ошибки или nil

* *metrics*:
    * `metrics.set(key, val, <ttl>)` установить значение метрики key, val может быть string, number. ttl по дефолту 300 секунд
    * `metrics.set_speed(key, val, <ttl>)` тоже самое, но считает скорость измерения
    * `metrics.set_counter_speed(key, val, <ttl>)` тоже самое, но считает только положительную скорость измерения
    * `metrics.get(key)` получить значение метрики key
    * `metrics.list()` список `{key:"xxx", value:"xxx", at:xxx}`
    * `metrics.delete(key)` удалить значение метрики key

* *postgres*:
    * `db, err = postgres.open({database="xxx", host="127.0.0.1", user="xxx", password="xxx"})` открыть коннект
    * `rows, err, column_count, row_count = db:query()` выполнить запрос
    * `db:close()` закрыть коннект

* *tcp*:
    * `telnet, err = tcp.open("xxx:xxx")` открыть коннект
    * `err = telnet:write("xxxx")` записать в коннект
    * `telnet:close()` закрыть коннект

* *tac*:
    * `scanner = tac.open("filepath")` открыть файл
    * `scanner:line()` получить последнюю линию (closure), в случае отсутвия таковой вернеться nil
    * `scanner:close()` закрыть файл

* *ioutil*:
    * `ioutil.readfile(filename)` вернуть содержимое файла

* *crypto*:
    * `crypto.md5(string)` вернуть md5 от строки

* *cmd*:
    * `state, err = cmd.exec(string)` выполнить exec через shell, возвращает `{"code"=0, "stderr"="", "stdout"=""}`

* *goruntime*:
    * `goruntime.goos` порт runtime.GOOS
    * `goruntime.goarch` порт runtime.GOARCH

* *parser*:
    позволяет загрузить при помощи https://golang.org/pkg/plugin/ библиотеку с переменная которая вернет интерфейс: `type Parser interface { ProcessData(string) (map[string]string, error) }`
    * `p, err = parser.load(filename.so, variable_name="NewParser")` загрузить плагин с экспортированным именем `variable_name` в filename.so
    * `table, err = p:parse(str)` обработать строчку

* *strings*:
    * `strings.split(str, delim)` порт golang strings.split()
    * `strings.has_prefix(str1, str2)` порт golang strings.hasprefix()
    * `strings.has_suffix(str1, str2)` порт golang strings.hassuffix()
    * `strings.trim(str1, str2)` порт golang strings.trim()

* *json*:
    * `json.encode(N)` lua-table в string
    * `json.decode(N)` string в lua-table

* *yaml*:
    * `yaml.decode(N)` string в lua-table

* *filepath*:
    * `filepath.base(filename)` порт golang filepath.Base()
    * `filepath.dir(filename)` порт golang filepath.Dir()
    * `filepath.ext(filename)` порт golang filepath.Ext()
    * `filepath.glob(mask)` порт golang filepath.Glob(), в случае ошибки возращает nil.

* *xmlpath*:
    * `table, err = xmlpath.parse(data, path)` возвращает таблицу с обработаной `data` по xmlpath `path`

* *http*:
    * `result, err = http.get(url)` возвращает `result = {body, code}` и ошибку, захардкожен 10секундный таймаут.
    * `result, err = http.unescape(url)` порт  url.QueryUnescape(query)
    * `result = http.escape(url)` порт url.QueryEscape(query)

* *goos*:
    * `stat = goos.stat(filename)` goos.stat возвращает таблицу с полями `stat = {size, is_dir, mod_time}`, в случае ошибки возращает nil.
    * `goos.pagesize()` возвращет pagesize

* *time*:
    * `time.sleep(N)` проспать N секунд
    * `time.unix()` время в секундах
    * `time.unix_nano()` время в наносекундах
    * `time.parse("2006-Jan-02", "2018-Mar-02")` golang порт time.Parse(), возвращает unixts и ошибку

* *log*:
    * `log.error(msg)` сообщение в лог с уровнем error
    * `log.info(msg)` с уровнем info

## Примеры плагинов

### Diskstat

![await](/img/await.png)

* пытается сопоставить блочному девайсу /mount/pount
* рассчитывает await и utilization по тем блочным девайсам, по которым ядро не ведет статистику (mdraid)

### IO

![io](/img/io-syscall.png)

* суммирует /proc/*pid*/io
* рассчитывает эффективность чтения из vfs cache как соотношение логического и физического чтения rchar/read_bytes

### SNMP

данные из /proc/net/snmp

![tcp connections](/img/tcp-speed.png)
![tcp error](/img/tcp-errors.png)

### CPU

cpu time (не нормированное по кол-ву CPU)
![process state](/img/cpu-time.png)

состояние процессов
![process state](/img/cpu-proc.png)

interrupts и context switching
![process state](/img/cpu-intr.png)

### NUMA

статистика выделения памяти и доступа к ней

![numa access](/img/numa-access.png)
![numa allocate](/img/numa-alloc.png)
