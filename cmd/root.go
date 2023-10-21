package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "wow-cli",
		Short: "A CLI for wow classic arena leaderboard app",
	}

	rootCmd.AddCommand(
		serveCmd(),
		migrateCmd(),
		fetchCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("Error executing root command: %v", err)
		os.Exit(1)
	}
}
