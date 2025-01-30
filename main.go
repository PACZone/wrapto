package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PACZone/wrapto/core"
	logger "github.com/PACZone/wrapto/log"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) {
	logger.Info("running wrapto", "version", StringVersion())
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	c, err := core.NewCore(ctx, cancel, args[0])
	exitOnError(cmd, err)

	go c.Start()

	<-ctx.Done()
	<-time.After(time.Second * 5)
	logger.Info("shutdown")
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "wrapto",
		Version: StringVersion(),
		Run:     run,
	}

	err := rootCmd.Execute()
	exitOnError(rootCmd, err)
}

func exitOnError(cmd *cobra.Command, err error) {
	if err != nil {
		cmd.PrintErrln(err.Error())
		os.Exit(1)
	}
}
