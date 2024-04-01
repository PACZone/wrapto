package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/PacmanHQ/teleport"
	"github.com/PacmanHQ/teleport/core"
	cobra "github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "teleport-cli",
		Version: teleport.StringVersion(),
	}

	runCommand(rootCmd)

	err := rootCmd.Execute()
	ExitOnError(rootCmd, err)
}

func ExitOnError(cmd *cobra.Command, err error) {
	if err != nil {
		cmd.PrintErr(err.Error())
		os.Exit(1)
	}
}

func runCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs a mainnet instance of Teleport", // TODO add a testnet mode for me.
	}
	parentCmd.AddCommand(run)

	psbn := run.Flags().Int("pactus-start-block", 0, "The start block number for Pactus listener")
	pson := run.Flags().Int("polygon-start-order", 0, "The start order number for Polygon listener")

	run.Run = func(cmd *cobra.Command, _ []string) {
		c, err := core.NewCore(*psbn, *pson)
		if err != nil {
			ExitOnError(cmd, err)
		}

		c.Start()
		
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan
	}
}
