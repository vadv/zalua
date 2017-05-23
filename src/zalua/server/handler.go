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
	log.Printf("[INFO] %s: accept\n", requestId)

	buf := make([]byte, settings.MaxSizeRequest())
	conn.SetReadDeadline(time.Now().Add(settings.TimeoutRead()))
	n, err := conn.Read(buf[:])
	if err != nil {
		log.Printf("[ERROR] %s: can't read request: %s\n", requestId, err.Error())
		return
	}
	request := string(buf[0:n])
	response := ""
	log.Printf("[INFO] %s: request: '%s'\n", requestId, request)

	defer func() {
		log.Printf("[INFO] %s: response: '%s'\n", requestId, response)
	}()

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
		sort.Strings(list)
		response = strings.Join(list, "\n")

	// list of running plugins
	case request == protocol.LIST_OF_PLUGINS:
		list := dsl.ListOfPlugins()
		sort.Strings(list)
		response = strings.Join(list, "\n")

	// get value of metric
	case strings.HasPrefix(request, protocol.GET_METRIC_VALUE):
		data := strings.Split(request, protocol.GET_METRIC_VALUE)
		if len(data) == 2 {
			val, ok := storage.Box.Get(strings.Trim(data[1], " "))
			if ok {
				response = val
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
