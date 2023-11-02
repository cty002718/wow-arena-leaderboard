package dao

import (
	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"gorm.io/gorm"
)

type ILeaderboardDao interface {
	Create(*orm.Leaderboard) error
	GetLatest(seasonId int, bracket string) (*orm.Leaderboard, error)
}

type LeaderboardDao struct {
	DB *gorm.DB
}

func NewLeaderboardDao(db *gorm.DB) ILeaderboardDao {
	return &LeaderboardDao{
		DB: db,
	}
}

func (dao *LeaderboardDao) Create(leaderboard *orm.Leaderboard) error {
	return dao.DB.Create(leaderboard).Error
}

func (dao *LeaderboardDao) GetLatest(seasonId int, bracket string) (*orm.Leaderboard, error) {
	var leaderboard orm.Leaderboard
	err := dao.DB.Where("season_id = ? AND bracket = ?", seasonId, bracket).
		Order("created_at DESC").
		First(&leaderboard).Error

	if err != nil {
		return nil, err
	}

	return &leaderboard, nil
}
