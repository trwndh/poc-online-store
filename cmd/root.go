package cmd

import (
	"github.com/spf13/cobra"
	"github.com/trwndh/poc-online-store/pkg/logger"
	"go.uber.org/zap"
)

var RootCmd = &cobra.Command{
	Use:   "poc-online-store",
	Short: "poc-online-store service",
}

func Execute() {
	log := logger.Log
	RootCmd.AddCommand(StartHTTP)
	if err := RootCmd.Execute(); err != nil {
		log.Error("error when execute root command", zap.Error(err))
	}
}
