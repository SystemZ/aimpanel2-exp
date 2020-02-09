package exit

import (
	"gitlab.com/systemz/aimpanel2/master/events"
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

		events.SSE.Shutdown()
		EXIT = true

		os.Exit(1)
	}()
}
