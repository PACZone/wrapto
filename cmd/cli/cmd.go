package main

import (
	"os"

	"github.com/PacmanHQ/teleport"
	cobra "github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "teleport-cli",
		Version: teleport.StringVersion(),
	}

	RunCommand(rootCmd)

	err := rootCmd.Execute()
	ExitOnError(rootCmd, err)
}

func ExitOnError(cmd *cobra.Command, err error) {
	if err != nil {
		cmd.PrintErr(err.Error())
		os.Exit(1)
	}
}
