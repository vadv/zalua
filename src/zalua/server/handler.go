package server

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"zalua/dsl"
	"zalua/protocol"
	"zalua/settings"
	"zalua/storage"
)

var random = rand.New(rand.NewSource(time.Now().Unix()))

// генерируем уникальный номер
func requestId() int64 {
	return random.Int63n(100000000)
}

// обслуживание клиента
func ClientHandler(conn net.Conn) {

	defer conn.Close()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover fatal error: %v\n", r)
		}
	}()

	requestId := fmt.Sprintf("request-id-%d", requestId())

	buf := make([]byte, settings.MaxSizeRequest())
	conn.SetReadDeadline(time.Now().Add(settings.TimeoutRead()))
	n, err := conn.Read(buf[:])
	if err != nil {
		log.Printf("[ERROR] %s: can't read request: %s\n", requestId, err.Error())
		return
	}
	request := string(buf[0:n])
	if request == protocol.COMMAND_SERVER_ID {
		// чтобы не засорять логи, через health-check проверяем ServerId
		conn.Write([]byte(settings.ServerId()))
		return
	}
	response := ""

	switch {

	// ping-pong
	case request == protocol.PING:
		response = protocol.PONG

	// command-kill
	case request == protocol.COMMAND_KILL:
		log.Printf("[FATAL] kill server now!\n")
		os.Exit(100)

	// list of metrics
	case request == protocol.LIST_OF_METRICS:
		list := storage.Box.List()
		stringList := []string{}
		for _, item := range list {
			tagsArr := []string{}
			if len(item.GetTags()) > 0 {
				for k, v := range item.GetTags() {
					tagsArr = append(tagsArr, fmt.Sprintf("%s=%s", k, v))
				}
			}
			stringList = append(stringList, fmt.Sprintf("%s\t%s\t\t%s\t\t%d", item.GetMetric(), strings.Join(tagsArr, " "), item.GetValue(), item.GetCreatedAt()))
		}
		sort.Strings(stringList)
		response = strings.Join(stringList, "\n")

	// list of running plugins
	case request == protocol.LIST_OF_PLUGINS:
		list := dsl.ListOfPlugins()
		sort.Strings(list)
		response = strings.Join(list, "\n")

	// get value of metric
	case strings.HasPrefix(request, protocol.GET_METRIC_VALUE):
		data := strings.Split(request, " ")
		if len(data) >= 2 {
			metric := strings.Trim(data[1], " ")
			tags := make(map[string]string, 0)
			if len(data) > 2 {
				// key1=val2 key2=val2
				for _, str := range data[2:] {
					strData := strings.Split(str, "=")
					if len(strData) == 2 {
						tags[strData[0]] = strData[1]
					}
				}
			}
			val, ok := storage.Box.Get(metric, tags)
			if ok {
				response = val.GetValue()
			} else {
				response = protocol.UNKNOWN_METRIC
			}
		} else {
			response = protocol.COMMAND_ERROR
		}

	// unknown metric
	default:
		response = protocol.UNKNOWN_COMMAND
	}

	// чтобы не получить EOF при чтении
	if response == "" {
		response = protocol.EMPTY
	}

	// записываем ответ
	conn.Write([]byte(response))
}
