package dao

import (
	"time"

	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ISeasonDao interface {
	FindById(id int) (*orm.Season, error)
	FindOrCreate(id int, startedAt time.Time, endedAt time.Time) (*orm.Season, error)
}

type SeasonDao struct {
	DB *gorm.DB
}

func NewSeasonDao(db *gorm.DB) ISeasonDao {
	return &SeasonDao{
		DB: db,
	}
}

func (dao *SeasonDao) FindById(id int) (*orm.Season, error) {
	season := orm.Season{}

	err := dao.DB.Where(&orm.Season{
		ID: id,
	}).First(&season).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Failed to find season by id")
	}

	return &season, nil
}

func (dao *SeasonDao) FindOrCreate(id int, startedAt time.Time, endedAt time.Time) (*orm.Season, error) {
	season := orm.Season{}

	err := dao.DB.Where(&orm.Season{
		ID: id,
	}).Attrs(&orm.Season{
		StartedAt: startedAt,
		EndedAt:   endedAt,
	}).FirstOrCreate(&season).Error
	if err != nil {
		return nil, errors.Wrap(err, "Failed to find or create season")
	}

	return &season, nil
}
