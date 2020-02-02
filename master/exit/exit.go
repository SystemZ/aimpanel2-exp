package exit

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	EXIT bool = false
)

func CheckForExitSignal() {
	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		_ = <-sigc

		//Check how it works :)
		//events.SSE.Shutdown()

		EXIT = true
		os.Exit(1)
	}()
}
