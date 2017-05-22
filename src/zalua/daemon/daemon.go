package daemon

// https://github.com/golang/go/issues/227
// https://habrahabr.ru/post/187668/

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	envDaemonName  = "_runned_is_daemon"
	envDaemonValue = "1"
)

func Daemonize(pwd string) error {
	return reborn(pwd)
}

func reborn(workDir string) (err error) {
	if !IsDaemon() {
		var path string
		if path, err = filepath.Abs(os.Args[0]); err != nil {
			return
		}
		cmd := exec.Command(path) // daemonize without parametrs
		envVar := fmt.Sprintf("%s=%s", envDaemonName, envDaemonValue)
		cmd.Env = append(os.Environ(), envVar)
		if err = cmd.Start(); err != nil {
			return
		}
		// пытаемся продолжить как клиент
		time.Sleep(time.Second)
		return
	}
	if len(workDir) != 0 {
		if err = os.Chdir(workDir); err != nil {
			return
		}
	}
	return
}

func IsDaemon() bool {
	return os.Getenv(envDaemonName) == envDaemonValue
}
