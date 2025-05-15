package service

import (
	_itemManagingModel "github.com/lekchan000/isekai-shop-api/pkg/itemManaging/model"
	_itemShopModel "github.com/lekchan000/isekai-shop-api/pkg/itemShop/model"
)

type ItemManagingService interface {
	Creating(itemCreatingReq *_itemManagingModel.ItemCreatingReq) (*_itemShopModel.Item, error)
}
