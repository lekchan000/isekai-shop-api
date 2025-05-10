package repository

type itemShopRepository struct{}

func NewItemShopRepositoryImpl() ItemShopRepository {
	return &itemShopRepository{}
}
