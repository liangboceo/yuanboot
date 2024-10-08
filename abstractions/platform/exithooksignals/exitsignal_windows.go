package exithooksignals

import (
	"fmt"
	"github.com/liangboceo/yuanboot/abstractions"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func HookSignals(host abstractions.IServiceHost) {
	quitSig := make(chan os.Signal)
	signal.Notify(
		quitSig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)

	go func() {
		var sig os.Signal
		for {
			sig = <-quitSig
			fmt.Println()
			switch sig {
			// graceful stop
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT:
				host.StopApplicationNotify()
				host.Shutdown()
				// terminate now
			}
			time.Sleep(time.Second * 15)
		}
	}()
}
