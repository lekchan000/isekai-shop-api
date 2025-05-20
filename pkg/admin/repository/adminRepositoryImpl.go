package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/lekchan000/isekai-shop-api/database"
	"github.com/lekchan000/isekai-shop-api/entities"
	_adminException "github.com/lekchan000/isekai-shop-api/pkg/admin/exception"
)

type adminRepositoryImpl struct {
	db     database.Database
	logger echo.Logger
}

func NewAdminRepositoryImpl(db database.Database, logger echo.Logger) AdminRepository {
	return &adminRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *adminRepositoryImpl) Creating(adminEntity *entities.Admin) (*entities.Admin, error) {
	admin := new(entities.Admin)

	if err := r.db.Connect().Create(adminEntity).Scan(admin).Error; err != nil {
		r.logger.Errorf("Creating admin failed: %s", err.Error())
		return nil, &_adminException.AdminCreating{AdminID: adminEntity.ID}
	}

	return admin, nil
}

func (r *adminRepositoryImpl) FindByID(adminID string) (*entities.Admin, error) {
	admin := new(entities.Admin)

	if err := r.db.Connect().Where("id = ?", adminID).First(admin).Error; err != nil {
		r.logger.Errorf("Find admin by ID failed: %d", err.Error())
		return nil, &_adminException.AdminNotFound{AdminID: adminID}
	}

	return admin, nil
}
