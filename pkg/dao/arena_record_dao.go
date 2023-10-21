package dao

import (
	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"gorm.io/gorm"
)

type IArenaRecordDao interface {
	Create(arenaRecord *orm.ArenaRecord) error
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
