package main

import (
	"fmt"
	"log"
	"os"

	"zalua/daemon"
	"zalua/logger"
	"zalua/protocol"
	"zalua/server"
	"zalua/settings"
	"zalua/socket"
)

var BuildVersion = "unknown"

func main() {

	help := func() {
		fmt.Fprintf(os.Stderr, "%s commands:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "-v, -version, --version\n\tGet version\n")
		fmt.Fprintf(os.Stderr, "-k, -kill, --kill, kill \n\tKill server\n")
		fmt.Fprintf(os.Stderr, "-m, -metrics, --list-metrics, metrics\n\tList of known metrics\n")
		fmt.Fprintf(os.Stderr, "-p, -plugins, --plugins, plugins\n\tList of running plugins\n")
		fmt.Fprintf(os.Stderr, "-g, -get, --get, --get-metric, get <metric>\n\tGet metric value\n")
		fmt.Fprintf(os.Stderr, "-ping, --ping, ping\n\tPing pong game\n")
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-v", "--version", "-version":
			fmt.Printf("%s version: %s\n", os.Args[0], BuildVersion)
			os.Exit(1)
		case "-h", "-help", "--help":
			help()
		}
	}

	// демонизация
	if !socket.Alive() {
		if !daemon.IsDaemon() {
			fmt.Fprintf(os.Stderr, "Socket %s is not alive, start daemon\n", settings.SocketPath())
		}
		if err := daemon.Daemonize("/"); err != nil {
			fmt.Fprintf(os.Stderr, "Daemonize error: %s\n", err.Error())
			os.Exit(2)
		}
		if !daemon.IsDaemon() {
			fmt.Fprintf(os.Stderr, "Daemon starting, continue as client\n")
		}
	}

	if daemon.IsDaemon() {

		// настройка лога
		fd, err := logger.GetLogFD()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't open log file: %s\n", err.Error())
			os.Exit(3)
		}
		logger.DupFdToStd(fd)
		log.SetOutput(fd)
		log.Printf("[INFO] Start server\n")

		// мы должны запустить сервис
		server.DoInit()
		if err := socket.ListenLoop(server.ClientHandler); err != nil {
			log.Printf("[FATAL] Listen server: %s\n", err.Error())
		}
		// здесь мы по идее не должны не когда оказаться
		os.Exit(5)
	}

	client, err := socket.GetClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(7)
	}

	// отправляем сообщение серверу
	msg := ""
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	switch arg {

	case "", "-ping", "--ping", "ping":
		msg = protocol.PING

	case "-k", "-kill", "--kill", "kill":
		msg = protocol.COMMAND_KILL

	case "-p", "-plugins", "--plugins", "--list-plugins", "plugins":
		if len(os.Args) != 2 {
			help()
		}
		msg = protocol.LIST_OF_PLUGINS

	case "-m", "-metrics", "--metrics", "--list-metrics", "metrics":
		if len(os.Args) != 2 {
			help()
		}
		msg = protocol.LIST_OF_METRICS

	case "-g", "-get", "--get", "--get-metric", "get":
		if len(os.Args) != 3 {
			help()
		}
		msg = fmt.Sprintf("%s %s", protocol.GET_METRIC_VALUE, os.Args[2])

	default:
		fmt.Fprintf(os.Stderr, "unknown command\n")
		os.Exit(8)
	}

	result, err := client.SendMessage(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "send message error: %s\n", err.Error())
		os.Exit(10)
	}
	switch result {
	case protocol.UNKNOWN_METRIC:
		// fmt.Fprintf(os.Stderr, "unknown metric '%s'\n", os.Args[2])
		fmt.Fprintf(os.Stdout, "")
		os.Exit(11)
	case protocol.UNKNOWN_COMMAND:
		fmt.Fprintf(os.Stderr, "unknown command\n")
		os.Exit(12)
	case protocol.COMMAND_ERROR:
		fmt.Fprintf(os.Stderr, "command error\n")
		os.Exit(13)
	case protocol.EMPTY:
		fmt.Fprintf(os.Stdout, "")
	default:
		fmt.Fprintf(os.Stdout, "%s\n", result)
	}

}
