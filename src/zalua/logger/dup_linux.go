package logger

import (
	"os"
	"syscall"
)

func DupFdToStd(fd *os.File) {
	syscall.Dup2(int(fd.Fd()), 1)
	syscall.Dup2(int(fd.Fd()), 2)
}
