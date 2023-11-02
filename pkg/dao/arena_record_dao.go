package dao

import (
	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IArenaRecordDao interface {
	Create(arenaRecord *orm.ArenaRecord) error
	FindByLeaderboardId(leaderboardId uuid.UUID) ([]*orm.ArenaRecord, error)
	FindByCharacter(characterId int64, seasonId int, bracket string) ([]*orm.ArenaRecord, error)
}

type ArenaRecordDao struct {
	DB *gorm.DB
}

func NewArenaRecordDao(db *gorm.DB) IArenaRecordDao {
	return &ArenaRecordDao{
		DB: db,
	}
}

func (dao *ArenaRecordDao) Create(arenaRecord *orm.ArenaRecord) error {
	return dao.DB.Create(arenaRecord).Error
}

func (dao *ArenaRecordDao) FindByLeaderboardId(leaderboardId uuid.UUID) ([]*orm.ArenaRecord, error) {
	var arenaRecords []*orm.ArenaRecord
	err := dao.DB.Preload("Character").Where("leaderboard_id = ?", leaderboardId).Order("rank ASC").Find(&arenaRecords).Error

	if err != nil {
		return nil, err
	}

	return arenaRecords, nil
}

func (dao *ArenaRecordDao) FindByCharacter(characterId int64, seasonId int, bracket string) ([]*orm.ArenaRecord, error) {
	var leaderboardIDs []uuid.UUID
	subQuery := dao.DB.Table("leaderboards").Select("id").Where("season_id = ? AND bracket = ?", seasonId, bracket)

	err := subQuery.Find(&leaderboardIDs).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"seasonId": seasonId,
			"bracket":  bracket,
		}).Error("failed to find leaderboard ids")
		return nil, err
	}

	var arenaRecords []*orm.ArenaRecord
	err = dao.DB.Preload("Leaderboard").Where("character_id = ? AND leaderboard_id IN ?", characterId, leaderboardIDs).Find(&arenaRecords).Error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"characterId":   characterId,
			"leaderboardId": leaderboardIDs,
		}).Error("failed to find arena records")
		return nil, err
	}

	return arenaRecords, nil
}
