package logger

import (
	"os"

	"zalua/settings"
)

var logFileFd *os.File

func GetLogFD() (*os.File, error) {
	if logFileFd != nil {
		return logFileFd, nil
	}
	fd, err := os.OpenFile(settings.LogPath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	logFileFd = fd
	return logFileFd, nil
}
