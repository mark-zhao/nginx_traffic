package common

import (
	"context"
	"glog"
	"os"
	"os/signal"
)

func SetSignal(cancelFunc context.CancelFunc, handler func()) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt)//, os.Interrupt, os.Kill )
	glog.Info("set signal success")
	sinalId := <-interrupt
	glog.Info("receive signal for exit; signal: ", sinalId)
	cancelFunc()
	handler()
}
