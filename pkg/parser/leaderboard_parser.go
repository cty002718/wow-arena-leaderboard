package parser

import (
	"time"

	"github.com/cty002718/wow-arena-leaderboard/pkg/apiclient"
	"github.com/cty002718/wow-arena-leaderboard/pkg/dao"
	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ILeaderboardParser interface {
	ParseLeaderboard(*apiclient.LeaderboardResponse) error
}

type LeaderboardParser struct {
	GameApiClient  apiclient.IGameApiClient
	SeasonDao      dao.ISeasonDao
	LeaderboardDao dao.ILeaderboardDao
	ServerDao      dao.IServerDao
	CharacterDao   dao.ICharacterDao
	ArenaRecordDao dao.IArenaRecordDao
}

func NewLeaderboardParser(
	gameApiClient apiclient.IGameApiClient,
	seasonDao dao.ISeasonDao,
	leaderboardDao dao.ILeaderboardDao,
	serverDao dao.IServerDao,
	characterDao dao.ICharacterDao,
	arenaRecordDao dao.IArenaRecordDao,
) ILeaderboardParser {
	return &LeaderboardParser{
		GameApiClient:  gameApiClient,
		SeasonDao:      seasonDao,
		LeaderboardDao: leaderboardDao,
		ServerDao:      serverDao,
		CharacterDao:   characterDao,
		ArenaRecordDao: arenaRecordDao,
	}
}

func (p *LeaderboardParser) ParseLeaderboard(leaderboard *apiclient.LeaderboardResponse) error {
	err := p.createSeasonIfNotExists(leaderboard.Season.ID)
	if err != nil {
		return errors.Wrap(err, "Failed to create season if not exists")
	}

	leaderboardModel := orm.Leaderboard{
		SeasonID: leaderboard.Season.ID,
		Bracket:  leaderboard.Bracket,
	}
	err = p.LeaderboardDao.Create(&leaderboardModel)
	if err != nil {
		return errors.Wrap(err, "Failed to create leaderboard")
	}

	for _, entry := range leaderboard.Entries {
		err := p.createOrUpdateCharacter(entry.Character, entry.Faction.Type)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"character_id":   entry.Character.ID,
				"character_name": entry.Character.Name,
				"server_slug":    entry.Character.Server.Slug,
				"error":          err,
			}).Errorf("Failed to create or update characte")
			continue
		}

		arenaRecordModel := orm.ArenaRecord{
			LeaderboardID: leaderboardModel.ID,
			CharacterID:   entry.Character.ID,
			Rank:          entry.Rank,
			Rating:        entry.Rating,
			Won:           entry.Stat.Won,
			Lost:          entry.Stat.Lost,
		}
		err = p.ArenaRecordDao.Create(&arenaRecordModel)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"leaderboard_id": leaderboardModel.ID,
				"character_id":   entry.Character.ID,
				"rank":           entry.Rank,
				"rating":         entry.Rating,
				"won":            entry.Stat.Won,
				"lost":           entry.Stat.Lost,
				"error":          err,
			}).Errorf("Failed to create arena record")
			continue
		}
	}

	logrus.WithFields(logrus.Fields{
		"season_id": leaderboard.Season.ID,
		"bracket":   leaderboard.Bracket,
	}).Infof("Leaderboard parsed successfully")

	return nil
}

func (p *LeaderboardParser) createSeasonIfNotExists(seasonId int) error {
	season, err := p.SeasonDao.FindById(seasonId)
	if err != nil {
		return errors.Wrap(err, "Failed to find season by id")
	}

	if season == nil {
		logrus.Infof("Season %d not found, fetching from api...", seasonId)
		seasonResponse, err := p.GameApiClient.GetSeason(seasonId)
		if err != nil {
			return errors.Wrap(err, "Failed to fetch season from api")
		}

		_, err = p.SeasonDao.FindOrCreate(
			seasonResponse.ID,
			time.Unix(seasonResponse.StartedAt/1000, 0),
			time.Unix(seasonResponse.EndedAt/1000, 0),
		)
		if err != nil {
			return errors.Wrap(err, "Failed to find or create season")
		}
		logrus.Infof("Season %d created", seasonId)
	}
	return nil
}

func (p *LeaderboardParser) createOrUpdateCharacter(character apiclient.Character, faction string) error {
	err := p.createServerIfNotExists(character.Server.Slug)
	if err != nil {
		return errors.Wrap(err, "Failed to create server if not exists")
	}

	characterModel := orm.Character{
		ID:       character.ID,
		Name:     character.Name,
		ServerID: character.Server.ID,
		Faction:  faction,
	}

	err = p.CharacterDao.CreateOrUpdate(&characterModel)
	if err != nil {
		return errors.Wrap(err, "Failed to create or update character")
	}

	return nil
}

func (p *LeaderboardParser) createServerIfNotExists(serverSlug string) error {
	server, err := p.ServerDao.FindBySlug(serverSlug)
	if err != nil {
		return errors.Wrap(err, "Failed to find server by slug")
	}

	if server == nil {
		logrus.Infof("Server %s not found, fetching from api...", serverSlug)
		serverResponse, err := p.GameApiClient.GetServerInfo(serverSlug)
		if err != nil {
			return errors.Wrap(err, "Failed to fetch server from api")
		}

		_, err = p.ServerDao.FindOrCreate(
			serverResponse.ID,
			serverResponse.Name,
			serverResponse.Slug,
			serverResponse.ServerType.Type,
		)
		if err != nil {
			return errors.Wrap(err, "Failed to find or create server")
		}
	}

	return nil
}
