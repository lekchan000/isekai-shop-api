package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/lekchan000/isekai-shop-api/database"
	"github.com/lekchan000/isekai-shop-api/entities"

	_itemManagingException "github.com/lekchan000/isekai-shop-api/pkg/itemManaging/exception"
	_itemManagingModel "github.com/lekchan000/isekai-shop-api/pkg/itemManaging/model"
)

type itemManagingRepositoryImpl struct {
	db     database.Database
	logger echo.Logger
}

func NewItemManagingRepositoryImpl(db database.Database, logger echo.Logger) *itemManagingRepositoryImpl {
	return &itemManagingRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *itemManagingRepositoryImpl) Creating(itemEntity *entities.Item) (*entities.Item, error) {
	item := new(entities.Item)

	if err := r.db.Connect().Create(itemEntity).Scan(item).Error; err != nil {
		r.logger.Errorf("Creating item failed: %s", err.Error())
		return nil, &_itemManagingException.ItemCreating{}
	}

	return item, nil
}

func (r *itemManagingRepositoryImpl) Editing(itemID uint64, itemEditingReq *_itemManagingModel.ItemEditingReq) (uint64, error) {
	if err := r.db.Connect().Model(&entities.Item{}).Where("id = ?", itemID).Updates(itemEditingReq).Error; err != nil {
		r.logger.Errorf("Editing Item Failed : %s", err.Error())
		return 0, &_itemManagingException.ItemEditing{}
	}
	return itemID, nil
}

func (r *itemManagingRepositoryImpl) Archiving(itemID uint64) error {
	if err := r.db.Connect().Table("items").Where("id = ?", itemID).Update("is_archive", true).Error; err != nil {
		r.logger.Errorf("Archiving item failed : %s", err.Error())
		return &_itemManagingException.ItemArchiving{ItemID: itemID}
	}
	return nil
}
