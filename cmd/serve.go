package cmd

import (
	"github.com/cty002718/wow-arena-leaderboard/pkg/controller"
	"github.com/gin-gonic/gin"
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
	logrus.Info("Starting server...")
	router := gin.Default()
	setupPublicRoute(router)
	router.Run(":8080")
}

func setupPublicRoute(router *gin.Engine) {
	publicRouteV1 := router.Group("public/api/v1")
	publicRouteV1.GET("/leaderboard", controller.GetLatestLeaderboard)
	publicRouteV1.GET("/character/:character_id/:season/:bracket", controller.GetCharacterArenaRecord)
}
