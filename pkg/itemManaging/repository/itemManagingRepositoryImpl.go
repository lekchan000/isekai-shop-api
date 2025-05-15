package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/lekchan000/isekai-shop-api/entities"
	"gorm.io/gorm"

	_itemManagingException "github.com/lekchan000/isekai-shop-api/pkg/itemManaging/exception"
)

type itemManagingRepositoryImpl struct {
	db     *gorm.DB
	logger echo.Logger
}

func NewItemManagingRepositoryImpl(db *gorm.DB, logger echo.Logger) *itemManagingRepositoryImpl {
	return &itemManagingRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *itemManagingRepositoryImpl) Creating(itemEntity *entities.Item) (*entities.Item, error) {
	item := new(entities.Item)

	if err := r.db.Create(itemEntity).Scan(item).Error; err != nil {
		r.logger.Errorf("Creating item failed: %s", err.Error())
		return nil, &_itemManagingException.ItemCreating{}
	}

	return item, nil
}
