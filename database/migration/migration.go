package main

import (
	"github.com/lekchan000/isekai-shop-api/config"
	"github.com/lekchan000/isekai-shop-api/database"
	"github.com/lekchan000/isekai-shop-api/entities"
)

func main() {
	conf := config.ConfigGetting()
	db := database.NewPostgresDatabase(conf.Database)

	tx := db.Connect().Begin()

	playerMigration(db)
	adminMigration(db)
	itemMigration(db)
	playerCoinMigration(db)
	inventoryMigration(db)
	purchaseHistoryMigration(db)

	tx.Commit()
	if tx.Error != nil {
		tx.Rollback()
		panic(tx.Error)
	}
}

func playerMigration(db database.Database) {
	db.Connect().Migrator().CreateTable(&entities.Player{})
}

func adminMigration(db database.Database) {
	db.Connect().Migrator().CreateTable(&entities.Admin{})
}

func itemMigration(db database.Database) {
	db.Connect().Migrator().CreateTable(&entities.Item{})
}

func playerCoinMigration(db database.Database) {
	db.Connect().Migrator().CreateTable(&entities.PlayerCoin{})
}

func inventoryMigration(db database.Database) {
	db.Connect().Migrator().CreateTable(&entities.Inventory{})
}

func purchaseHistoryMigration(db database.Database) {
	db.Connect().Migrator().CreateTable(&entities.PurchaseHistory{})
}
