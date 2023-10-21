package dao

import (
	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IServerDao interface {
	FindById(id int) (*orm.Server, error)
	FindBySlug(slug string) (*orm.Server, error)
	FindOrCreate(id int, name string, slug string, serverType string) (*orm.Server, error)
}

type ServerDao struct {
	DB *gorm.DB
}

func NewServerDao(db *gorm.DB) IServerDao {
	return &ServerDao{
		DB: db,
	}
}

func (dao *ServerDao) FindById(id int) (*orm.Server, error) {
	server := orm.Server{}
	err := dao.DB.Where(&orm.Server{
		ID: id,
	}).First(&server).Error
	if err != nil {
		return nil, errors.Wrap(err, "Failed to find server by id")
	}

	return &server, nil
}

func (dao *ServerDao) FindBySlug(slug string) (*orm.Server, error) {
	server := orm.Server{}
	err := dao.DB.Where(&orm.Server{
		Slug: slug,
	}).First(&server).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Failed to find server by slug")
	}

	return &server, nil
}

func (dao *ServerDao) FindOrCreate(id int, name string, slug string, serverType string) (*orm.Server, error) {
	server := orm.Server{}

	err := dao.DB.Where(&orm.Server{
		ID: id,
	}).Attrs(&orm.Server{
		Name: name,
		Slug: slug,
		Type: serverType,
	}).FirstOrCreate(&server).Error
	if err != nil {
		return nil, errors.Wrap(err, "Failed to find or create server")
	}

	return &server, nil
}
