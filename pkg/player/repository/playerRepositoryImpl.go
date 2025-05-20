package repository

import (
	"github.com/labstack/echo/v4"
	"github.com/lekchan000/isekai-shop-api/database"
	"github.com/lekchan000/isekai-shop-api/entities"

	_playerException "github.com/lekchan000/isekai-shop-api/pkg/player/exception"
)

type playerRepositoryImpl struct {
	db     database.Database
	logger echo.Logger
}

func NewPlayerRepositoryImpl(db database.Database, logger echo.Logger) PlayerRepository {
	return &playerRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *playerRepositoryImpl) Creating(playerEntity *entities.Player) (*entities.Player, error) {
	player := new(entities.Player)

	if err := r.db.Connect().Create(playerEntity).Scan(player).Error; err != nil {
		r.logger.Errorf("Creating player failed: %s", err.Error())
		return nil, &_playerException.PlayerCreating{PlayerID: playerEntity.ID}
	}

	return player, nil
}

func (r *playerRepositoryImpl) FindByID(playerID string) (*entities.Player, error) {

	player := new(entities.Player)

	if err := r.db.Connect().Where("id = ?", playerID).First(player).Error; err != nil {
		r.logger.Errorf("Find player by ID failed: %d", err.Error())
		return nil, &_playerException.PlayerNotFound{PlayerID: playerID}
	}

	return player, nil
}
