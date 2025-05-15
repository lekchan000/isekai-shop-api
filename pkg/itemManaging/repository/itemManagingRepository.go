package repository

import "github.com/lekchan000/isekai-shop-api/entities"

type ItemManagingRepository interface {
	Creating(itemEntity *entities.Item) (*entities.Item, error)
}
