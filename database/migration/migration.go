package main

import (
	"fmt"

	"github.com/lekchan000/isekai-shop-api/config"
	"github.com/lekchan000/isekai-shop-api/database"
	"github.com/lekchan000/isekai-shop-api/entities"
)

func main() {
	conf := config.ConfigGetting()
	db := database.NewPostgresDatabase(conf.Database)

	// Run migrations
	playerMigration(db)
	adminMigration(db)
	itemMigration(db)
	playerCoinMigration(db)
	inventoryMigration(db)
	purchaseHistoryMigration(db)
	fmt.Println("Migrations successfully.")
}

func playerMigration(db database.Database) {
	db.ConnectionGetting().Migrator().CreateTable(&entities.Player{})
}

func adminMigration(db database.Database) {
	db.ConnectionGetting().Migrator().CreateTable(&entities.Admin{})
}

func itemMigration(db database.Database) {
	db.ConnectionGetting().Migrator().CreateTable(&entities.Item{})
}

func playerCoinMigration(db database.Database) {
	db.ConnectionGetting().Migrator().CreateTable(&entities.PlayerCoin{})
}

func inventoryMigration(db database.Database) {
	db.ConnectionGetting().Migrator().CreateTable(&entities.Inventory{})
}

func purchaseHistoryMigration(db database.Database) {
	db.ConnectionGetting().Migrator().CreateTable(&entities.PurchaseHistory{})
}
