package service

import (
	_ItemShopModel "github.com/lekchan000/isekai-shop-api/pkg/itemShop/model"
)

type ItemShopService interface {
	Listing() ([]*_ItemShopModel.Item, error)
}
