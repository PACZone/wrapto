package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/PacmanHQ/teleport/core"
)

func main() {
	c := core.NewCore()
	c.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigChan
}
