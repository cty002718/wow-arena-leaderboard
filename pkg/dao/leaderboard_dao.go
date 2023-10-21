package dao

import (
	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"gorm.io/gorm"
)

type ILeaderboardDao interface {
	Create(*orm.Leaderboard) error
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
