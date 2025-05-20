package repository

import "github.com/lekchan000/isekai-shop-api/entities"

type AdminRepository interface {
	Creating(adminEntity *entities.Admin) (*entities.Admin, error)
	FindByID(adminID string) (*entities.Admin, error)
}
