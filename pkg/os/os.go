package os

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitSignalChina() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
