package cmd

import (
	"os"
	"time"

	"github.com/cty002718/wow-arena-leaderboard/pkg/apiclient"
	"github.com/cty002718/wow-arena-leaderboard/pkg/dao"
	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"github.com/cty002718/wow-arena-leaderboard/pkg/parser"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func fetchCmd() *cobra.Command {
	fetchCmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch and parse data from tw.api.blizzard.com",
	}
	fetchCmd.AddCommand(
		leaderboardFetchCmd(),
		priceFetchCmd(),
	)

	return fetchCmd
}

func leaderboardFetchCmd() *cobra.Command {
	leaderboardFetchCmd := &cobra.Command{
		Use:   "leaderboard",
		Short: "Fetch and parse leaderboard data from tw.api.blizzard.com",
		Run: func(cmd *cobra.Command, args []string) {
			all, _ := cmd.Flags().GetBool("all")
			leaderboardFetch(all)
		},
	}
	leaderboardFetchCmd.Flags().BoolP("all", "a", false, "Fetch all seasons")

	return leaderboardFetchCmd
}

func priceFetchCmd() *cobra.Command {
	priceFetchCmd := &cobra.Command{
		Use:   "price",
		Short: "Fetch and parse price data from tw.api.blizzard.com",
		Run: func(cmd *cobra.Command, args []string) {
			priceFetch()
		},
	}

	return priceFetchCmd
}

type LeaderboardFetcher struct {
	LeaderboardParser parser.ILeaderboardParser
	GameApiClient     apiclient.IGameApiClient
}

func NewLeaderboardFetcher(LeaderboardParser parser.ILeaderboardParser, GameApiClient apiclient.IGameApiClient) *LeaderboardFetcher {
	return &LeaderboardFetcher{
		LeaderboardParser: LeaderboardParser,
		GameApiClient:     GameApiClient,
	}
}

func leaderboardFetch(all bool) {
	TokenApiClient := apiclient.NewTokenApiClient()
	accessToken, err := TokenApiClient.GetClientAccessToken()
	if err != nil {
		logrus.WithError(err).Error("Failed to get client access token")
		return
	}

	GameApiClient := apiclient.NewGameApiClient(accessToken)

	db, err := orm.OpenDB()
	if err != nil {
		logrus.WithError(err).Error("Failed to open db")
		return
	}

	seasonDao := dao.NewSeasonDao(db)
	leaderboardDao := dao.NewLeaderboardDao(db)
	serverDao := dao.NewServerDao(db)
	characterDao := dao.NewCharacterDao(db)
	arenaRecordDao := dao.NewArenaRecordDao(db)
	LeaderboardParser := parser.NewLeaderboardParser(
		GameApiClient,
		seasonDao,
		leaderboardDao,
		serverDao,
		characterDao,
		arenaRecordDao,
	)

	Fetcher := NewLeaderboardFetcher(LeaderboardParser, GameApiClient)

	if all {
		seasons := []int{5, 6, 7, 8}
		brackets := []string{"2v2", "3v3", "5v5"}
		for _, seasonId := range seasons {
			for _, bracket := range brackets {
				Fetcher.fetch(seasonId, bracket)
				time.Sleep(2 * time.Second)
			}
		}
	} else {
		seasonId := cast.ToInt(os.Getenv("CURRENT_SEASON_ID"))
		brackets := []string{"2v2", "3v3", "5v5"}
		for _, bracket := range brackets {
			Fetcher.fetch(seasonId, bracket)
		}
	}
}

func (l *LeaderboardFetcher) fetch(seasonId int, bracket string) {
	logrus.WithFields(logrus.Fields{
		"season_id": seasonId,
		"bracket":   bracket,
	}).Infof("Fetching leaderboard...")
	leaderboardResponse, err := l.GameApiClient.GetLeaderboard(seasonId, bracket)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"season_id": seasonId,
			"bracket":   bracket,
		}).WithError(err).Warnf("Failed to fetch leaderboard")
		return
	}

	logrus.WithFields(logrus.Fields{
		"season_id": seasonId,
		"bracket":   bracket,
	}).Infof("Saving leaderboard into DB...")
	err = l.LeaderboardParser.ParseLeaderboard(leaderboardResponse)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"season_id": seasonId,
			"bracket":   bracket,
		}).WithError(err).Warnf("Failed to parse leaderboard")
	}
}

func priceFetch() {
	TokenApiClient := apiclient.NewTokenApiClient()
	accessToken, err := TokenApiClient.GetClientAccessToken()
	if err != nil {
		logrus.WithError(err).Error("Failed to get client access token")
		return
	}

	GameApiClient := apiclient.NewGameApiClient(accessToken)
	db, err := orm.OpenDB()
	if err != nil {
		logrus.WithError(err).Error("Failed to open db")
		return
	}
	TokenPriceLogDao := dao.NewTokenPriceLogDao(db)

	logrus.Info("Fetching token price...")
	tokenPriceResponse, err := GameApiClient.GetTokenPrice()
	if err != nil {
		logrus.WithError(err).Error("Failed to fetch token price")
		return
	}

	logrus.Info("Saving token price into DB...")
	timeInSeconds := tokenPriceResponse.LastUpdatedTimestamp / 1000
	lastUpdatedTime := time.Unix(timeInSeconds, 0)

	tokenPriceLog, err := TokenPriceLogDao.FindOrCreate(lastUpdatedTime, tokenPriceResponse.Price)
	if err != nil {
		logrus.WithError(err).Error("Failed to save token price into DB")
		return
	}

	logrus.WithFields(map[string]interface{}{
		"last_updated_time": tokenPriceLog.LastUpdatedTime.Local(),
		"price":             tokenPriceLog.Price,
	}).Info("Token price saved into DB")
}
