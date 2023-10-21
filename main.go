package main

import (
	"github.com/cty002718/wow-arena-leaderboard/cmd"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Warn("Failed to load .env file")
	}
}

func main() {
	cmd.Execute()
}
