package dao

import (
	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ICharacterDao interface {
	CreateOrUpdate(*orm.Character) error
}

type CharacterDao struct {
	DB *gorm.DB
}

func NewCharacterDao(db *gorm.DB) ICharacterDao {
	return &CharacterDao{
		DB: db,
	}
}

func (dao *CharacterDao) CreateOrUpdate(character *orm.Character) error {
	err := dao.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(character).Error
	if err != nil {
		return errors.Wrap(err, "Failed to create or update character")
	}
	return nil
}
