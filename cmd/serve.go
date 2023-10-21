package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func serveCmd() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the API",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}

	return serveCmd
}

func serve() {
	logrus.Info("Serving API")
}
