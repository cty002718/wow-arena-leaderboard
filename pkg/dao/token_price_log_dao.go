package dao

import (
	"time"

	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ITokenPriceLogDao interface {
	FindOrCreate(lastUpdateTime time.Time, price int) (*orm.TokenPriceLog, error)
}

type TokenPriceLogDao struct {
	DB *gorm.DB
}

func NewTokenPriceLogDao(db *gorm.DB) ITokenPriceLogDao {
	return &TokenPriceLogDao{
		DB: db,
	}
}

func (dao *TokenPriceLogDao) FindOrCreate(lastUpdateTime time.Time, price int) (*orm.TokenPriceLog, error) {
	tokenPriceLog := orm.TokenPriceLog{}

	err := dao.DB.Where(&orm.TokenPriceLog{
		LastUpdatedTime: lastUpdateTime,
	}).Attrs(&orm.TokenPriceLog{
		Price: price,
	}).FirstOrCreate(&tokenPriceLog).Error
	if err != nil {
		return nil, errors.Wrap(err, "Failed to find or create token price log")
	}

	return &tokenPriceLog, nil
}
