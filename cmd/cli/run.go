package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/PacmanHQ/teleport/core"
	"github.com/spf13/cobra"
)

func RunCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs a mainnet instance of Teleport", // TODO add a testnet mode for me.
	}
	parentCmd.AddCommand(run)

	psbn := run.Flags().Int("pactus-start-block", 0, "The start block number for Pactus listener")
	pson := run.Flags().Int("polygon-start-order", 0, "The start order number for Polygon listener")

	run.Run = func(cmd *cobra.Command, _ []string) {
		c := core.NewCore(*psbn, *pson)
		c.Start()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan
	}
}
