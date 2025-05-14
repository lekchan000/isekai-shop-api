package repository

import (
	"github.com/lekchan000/isekai-shop-api/entities"
	_itemShopModel "github.com/lekchan000/isekai-shop-api/pkg/itemShop/model"
)

type ItemShopRepository interface {
	Listing(itemFilter *_itemShopModel.ItemFilter) ([]*entities.Item, error)
}
