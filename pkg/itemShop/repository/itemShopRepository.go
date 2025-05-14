package repository

import "github.com/lekchan000/isekai-shop-api/entities"

type ItemShopRepository interface {
	Listing() ([]*entities.Item, error)
}
