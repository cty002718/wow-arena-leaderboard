package main

import (
	"github.com/cty002718/wow-arena-leaderboard/cmd"
	"github.com/cty002718/wow-arena-leaderboard/pkg/ctx"
	"github.com/cty002718/wow-arena-leaderboard/pkg/dao"
	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Warn("Failed to load .env file")
	}
}

func main() {
	logrus.Info("Starting application...")
	provideDependencies(ctx.Container)

	cmd.Execute()
}

func provideDependencies(container *dig.Container) {
	container.Provide(orm.OpenDB)
	container.Provide(dao.NewArenaRecordDao)
	container.Provide(dao.NewLeaderboardDao)
	container.Provide(dao.NewCharacterDao)
	container.Provide(dao.NewSeasonDao)
	container.Provide(dao.NewServerDao)
	container.Provide(dao.NewTokenPriceLogDao)
}
