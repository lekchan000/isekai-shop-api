package main

import (
	"fmt"

	"github.com/lekchan000/isekai-shop-api/config"
	"github.com/lekchan000/isekai-shop-api/database"
)

func main() {
	conf := config.ConfigGetting()
	db := database.NewPostgresDatabase(conf.Database)

	fmt.Println(db.ConnectionGetting())
}
